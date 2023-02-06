package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/client/mcs"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/common"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/internal"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/model"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
	"github.com/unknwon/com"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var log = logging.Logger("service")

// @Summary Summary information display
// @Produce  json
// @Success 200 {object} internal.Response{data=service.SummaryResp}
// @Failure 500 {object} internal.Response
// @Router /summary [get]
func Summary(c *gin.Context) {
	appG := internal.Gin{C: c}
	var summary SummaryResp

	countFileSource, err := model.CountFileSource()
	if err != nil {
		appG.Response(http.StatusInternalServerError, internal.ERROR_SUMMARY_FAIL, nil)
		return
	}
	summary.CidsCount = countFileSource

	deals, err := model.CountDealByMinerDeal()
	if err != nil {
		log.Errorf("Summary get CountDealByMinerDeal failed,error: %v", err)
		appG.Response(http.StatusInternalServerError, internal.ERROR_SUMMARY_FAIL, nil)
		return
	}
	summary.DealsCount = deals

	providers, err := model.CountProviderMinerDeal()
	if err != nil {
		log.Errorf("Summary get IpfsNodeCount failed,error: %v", err)
		appG.Response(http.StatusInternalServerError, internal.ERROR_SUMMARY_FAIL, nil)
		return
	}
	summary.Providers = providers

	ipfsNodes, err := model.IpfsNodeCount()
	if err != nil {
		log.Errorf("Summary get IpfsNodeCount failed,error: %v", err)
		appG.Response(http.StatusInternalServerError, internal.ERROR_SUMMARY_FAIL, nil)
		return
	}
	summary.IpfsNodes = ipfsNodes

	hotDataSize, err := model.HotDataSize()
	if err != nil {
		log.Errorf("Summary get HotDataSize failed,error: %v", err)
		appG.Response(http.StatusInternalServerError, internal.ERROR_SUMMARY_FAIL, nil)
		return
	}
	coldDataSize, err := model.ColdDataSize()
	if err != nil {
		log.Errorf("Summary get ColdDataSize failed,error: %v", err)
		appG.Response(http.StatusInternalServerError, internal.ERROR_SUMMARY_FAIL, nil)
		return
	}
	var total int64
	for _, hot := range hotDataSize {
		total += hot.FileSize * hot.Num
	}
	for _, cold := range coldDataSize {
		total += cold.FileSize * cold.Num
	}
	summary.DataStored = total

	height, err := common.NewLotusClient(10).GetCurrentHeight()
	if err != nil {
		log.Errorf("Summary get height failed,error: %v", err)
		appG.Response(http.StatusInternalServerError, internal.ERROR_SUMMARY_FAIL, nil)
		return
	}
	summary.Height = height
	appG.Response(http.StatusOK, internal.SUCCESS, summary)
}

// @Summary Get a list of file storage information
// @Produce  json
// @param	field_name	query	string	false	"data_cid/file_name"
// @param	page		query	int		false	"Page number, starting from 0 by default"
// @param	size		query	int		false	"By default, there are 20 lines."
// @Success 200 {object} internal.Response{data=service.FileSourcePager}
// @Failure 500 {object} internal.Response
// @Router /files [get]
func GetSourceList(c *gin.Context) {
	appG := internal.Gin{C: c}
	fieldName := com.StrTo(c.Query("field_name")).String()
	page, err := com.StrTo(c.Query("page")).Int64()
	if err != nil {
		page = 0
	}
	size, err := com.StrTo(c.Query("size")).Int64()
	if err != nil {
		size = 20
	}

	count, pageNum, err := model.CountFileSourceList(fieldName, size)
	if err != nil {
		appG.Response(http.StatusInternalServerError, internal.ERROR_FILE_LIST_FAIL, nil)
		return
	}

	sourceList, err := model.FileSourceList(fieldName, page, size)
	if err != nil {
		appG.Response(http.StatusInternalServerError, internal.ERROR_FILE_LIST_FAIL, nil)
		return
	}

	result := make([]FileSourceResp, 0)
	for _, fileSource := range sourceList {
		var fsr FileSourceResp
		fsr.FileName = fileSource.FileName
		fsr.DataCid = fileSource.DataCid
		fsr.FileSize = fileSource.FileSize
		fsr.Status = string(fileSource.RebuildStatus)
		ipfsUrls := make([]string, 0)
		for _, url := range fileSource.IpfsUrls {
			ipfsUrls = append(ipfsUrls, url.IpfsUrl)
		}
		fsr.IpfsUrls = ipfsUrls
		fsr.HotBackups = len(ipfsUrls)
		fsr.RebuildStatus = true

		providerStatus := make([]ProviderInfo, 0)
		for _, p := range fileSource.MinerIds {
			providerStatus = append(providerStatus, ProviderInfo{
				ProviderId: p.MinerId,
				Status:     p.Status,
			})
		}
		fsr.Providers = providerStatus
		fsr.ColdBackups = len(providerStatus)
		if len(providerStatus) == 0 {
			fsr.NotFoundProvider = internal.GetMsg(internal.ERROR_RETRIEVE_FAIL)
		}
		result = append(result, fsr)
	}
	appG.Response(http.StatusOK, internal.SUCCESS, FileSourcePager{
		Total:     count,
		PageCount: pageNum,
		Sources:   result,
	})
}

// @Summary Retrieve file storage from filecoin to ipfs based on cid
// @Produce  json
// @Param cid path string true "ID"
// @Success 200 {object} internal.Response
// @Failure 500 {object} internal.Response
// @Router /cid/{cid} [get]
func GetCid(c *gin.Context) {
	appG := internal.Gin{C: c}
	cid := com.StrTo(c.Param("cid")).String()
	//  1. query peerId by indexer node
	model.UpdateSourceFileStatus(cid, model.REBUILD_INDEXING)
	peerData := common.NewIndexerClient().SendHttpGet(common.GET_PEER_URL, cid)

	peerIds := make(map[string]string, 0)
	for _, data := range peerData {
		if string(data) == "no results for query" {
			continue
		}
		var indexData IndexData
		if err := json.Unmarshal(data, &indexData); err != nil {
			log.Errorf("change to json failed,error: %v", err)
			model.UpdateSourceFileStatus(cid, model.REBUILD_INDEXING_FAILED)
			continue
		}
		if len(indexData.MultihashResults) > 0 && len(indexData.MultihashResults[0].ProviderResults) > 0 {
			peerId := indexData.MultihashResults[0].ProviderResults[0].Provider.ID
			peerIds[peerId] = peerId
		}
	}
	if len(peerIds) == 0 {
		model.UpdateSourceFileStatus(cid, model.REBUILD_INDEXING_FAILED)
		appG.Response(http.StatusInternalServerError, internal.ERROR_RETRIEVE_FAIL, nil)
		return
	} else {
		model.UpdateSourceFileStatus(cid, model.REBUILD_RETRIEVING)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("catch panic error message: %v \n", err)
			}
		}()

		var successFlag bool
		lotusClient := common.NewLotusClient()
		for _, peerId := range peerIds {
			log.Infof("start process peerId: %s", peerId)
			// 2. query minerId
			mp, err := model.FindMinerPeer(peerId)
			if err != nil {
				log.Warnf("get minerpeer failed,peerId:%s,error: %v,continue check next peerId", peerId, err)
				continue
			}

			fileName := mp.MinerId + "-" + cid
			savePath := filepath.Join(model.LotusSetting.DownloadDir, fileName)
			// 3. retrieveData
			if err := lotusClient.RetrieveData(mp.MinerId, cid, savePath); err != nil {
				model.UpdateSourceFileStatus(cid, model.REBUILD_RETRIEVING_FAILED)
				continue
			}

			// 4. upload file to ipfs
			model.UpdateSourceFileStatus(cid, model.REBUILD_UPLOADING)
			stat, err := os.Stat(savePath)
			if err != nil {
				log.Errorf("not found savepath: %s,error: %s", savePath, err)
				return
			}

			fileIpfs := make([]model.FileIpfs, 0)
			if stat.IsDir() {
				urls := UploaderDir(savePath, model.UploaderSetting.IpfsUrls)
				for _, u := range urls {
					fileIpfs = append(fileIpfs, model.FileIpfs{
						DataCid: cid,
						IpfsUrl: u,
					})
				}
			} else {
				objectName := path.Join(time.Now().Format("2006-01-02"), fileName)
				if _, err = mcs.UploadFile(context.TODO(), "rebuilder", objectName, savePath); err != nil {
					log.Errorf("upload file to mcs bucket failed, error: %+v", err)
					model.UpdateSourceFileStatus(cid, model.REBUILD_UPLOADING_FAILED)
					return
				}
				fileUrl, err := mcs.GetFile(context.TODO(), "rebuilder", objectName)
				if err != nil {
					log.Errorf("get file from mcs bucket failed, error: %+v", err)
					model.UpdateSourceFileStatus(cid, model.REBUILD_UPLOADING_FAILED)
					return
				}
				fileIpfs = append(fileIpfs, model.FileIpfs{
					DataCid: cid,
					IpfsUrl: fileUrl,
				})
			}
			if len(fileIpfs) > 0 {
				if err = model.InsertFileIpfs(fileIpfs); err == nil {
					model.UpdateSourceFileStatus(cid, model.REBUILD_SUCCESS)
					successFlag = true
					os.RemoveAll(savePath)
				}
			}
			break
		}
		if !successFlag {
			model.UpdateSourceFileStatus(cid, model.REBUILD_FAILED)
		}
	}()
	appG.Response(http.StatusOK, internal.SUCCESS, map[string]interface{}{
		"msg": "Submitted for processing",
	})
}

// @Summary upload file
// @Description
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /upload [post]
func UploadFile(c *gin.Context) {
	appG := internal.Gin{C: c}
	file, err := c.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, internal.INVALID_PARAMS, internal.GetMsg(internal.INVALID_PARAMS))
		return
	}

	basePath := "./upload/"
	filename := basePath + filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		appG.Response(http.StatusInternalServerError, internal.ERROR_UPLOAD_FAIL, internal.GetMsg(internal.ERROR_UPLOAD_FAIL))
		return
	}
	appG.Response(http.StatusOK, internal.SUCCESS, "Uploaded successfully!")
}

// @Summary retrieve file
// @Description
// @accept application/json
// @Success 200 {object} internal.Response
// @Router /retrieve [post]
func Retrieve(c *gin.Context) {
	appG := internal.Gin{C: c}
	var retrieveReq RetrieveReq
	if err := c.ShouldBindJSON(&retrieveReq); err != nil {
		appG.Response(http.StatusInternalServerError, internal.ERROR_CHANGETO_JSON, nil)
	}
	if retrieveReq.DataCid == "" {
		appG.Response(http.StatusBadRequest, internal.INVALID_PARAMS, internal.GetMsg(internal.INVALID_PARAMS))
	}

	var sf model.SourceFile
	sf.DataCid = retrieveReq.DataCid
	sf.CreateAt = time.Now()
	model.InsertSourceFile(&sf)

	model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_INDEXING)
	peerData := common.NewIndexerClient().SendHttpGet(common.GET_PEER_URL, retrieveReq.DataCid)
	peerIds := make(map[string]string, 0)
	for _, data := range peerData {
		if string(data) == "no results for query" {
			continue
		}
		var indexData IndexData
		if err := json.Unmarshal(data, &indexData); err != nil {
			log.Errorf("change to json failed,error: %v", err)
			model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_INDEXING_FAILED)
			continue
		}
		if len(indexData.MultihashResults) > 0 && len(indexData.MultihashResults[0].ProviderResults) > 0 {
			peerId := indexData.MultihashResults[0].ProviderResults[0].Provider.ID
			peerIds[peerId] = peerId
		}
	}
	if len(peerIds) == 0 {
		model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_INDEXING_FAILED)
		appG.Response(http.StatusInternalServerError, internal.ERROR_RETRIEVE_FAIL, nil)
		return
	} else {
		model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_RETRIEVING)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("catch panic error message: %v \n", err)
			}
		}()

		var successFlag bool
		var stat os.FileInfo
		lotusClient := common.NewLotusClient()
		for _, peerId := range peerIds {
			log.Infof("start process peerId: %s", peerId)
			// 2. query minerId
			mp, err := model.FindMinerPeer(peerId)
			if err != nil {
				log.Warnf("get minerpeer failed,peerId:%s,error: %v,continue check next peerId", peerId, err)
				continue
			}

			var fm model.FileMiner
			fm.DataCid = retrieveReq.DataCid
			fm.MinerId = mp.MinerId
			fm.Status = "StorageDealActive"
			model.InsertFileMiner(&fm)

			fileName := mp.MinerId + "-" + retrieveReq.DataCid
			savePath := filepath.Join(model.LotusSetting.DownloadDir, fileName)
			// 3. retrieveData
			if err := lotusClient.RetrieveData(mp.MinerId, retrieveReq.DataCid, savePath); err != nil {
				model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_RETRIEVING_FAILED)
				continue
			}

			// 4. upload file to ipfs
			model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_UPLOADING)
			stat, err = os.Stat(savePath)
			if err != nil {
				log.Errorf("not found savepath: %s,error: %s", savePath, err)
				return
			}

			fileIpfs := make([]model.FileIpfs, 0)
			if stat.IsDir() {
				urls := UploaderDir(savePath, model.UploaderSetting.IpfsUrls)
				for _, u := range urls {
					fileIpfs = append(fileIpfs, model.FileIpfs{
						DataCid: retrieveReq.DataCid,
						IpfsUrl: u,
					})
				}
			} else {
				objectName := path.Join(time.Now().Format("2006-01-02"), fileName)
				if _, err = mcs.UploadFile(context.TODO(), "rebuilder", objectName, savePath); err != nil {
					log.Errorf("upload file to mcs bucket failed, error: %+v", err)
					model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_UPLOADING_FAILED)
					return
				}
				fileUrl, err := mcs.GetFile(context.TODO(), "rebuilder", objectName)
				if err != nil {
					log.Errorf("get file from mcs bucket failed, error: %+v", err)
					model.UpdateSourceFileStatus(retrieveReq.DataCid, model.REBUILD_UPLOADING_FAILED)
					return
				}
				fileIpfs = append(fileIpfs, model.FileIpfs{
					DataCid: retrieveReq.DataCid,
					IpfsUrl: fileUrl,
				})
			}
			if len(fileIpfs) > 0 {
				if err = model.InsertFileIpfs(fileIpfs); err == nil {
					successFlag = true
					var sf model.SourceFile
					sf.DataCid = retrieveReq.DataCid
					sf.FileSize = stat.Size()
					sf.FileName = retrieveReq.DataCid
					sf.RebuildStatus = model.REBUILD_SUCCESS
					model.InsertSourceFile(&sf)
					os.RemoveAll(savePath)
				}
			}
			break
		}
		if !successFlag {
			var sf model.SourceFile
			sf.DataCid = retrieveReq.DataCid
			sf.FileSize = stat.Size()
			sf.FileName = retrieveReq.DataCid
			sf.RebuildStatus = model.REBUILD_FAILED
			model.InsertSourceFile(&sf)
		}
	}()
	appG.Response(http.StatusOK, internal.SUCCESS, map[string]interface{}{
		"msg": "Submitted for processing",
	})
}

func WatchIpfsNodeData() {
	ticker := time.NewTicker(8 * time.Minute)
	for {
		select {
		case <-ticker.C:
			fileIpfsList := model.FindFileIpfsList()
			for _, fileIpfs := range fileIpfsList {
				split := strings.Split(fileIpfs.IpfsUrl, "/ipfs/")
				if alive := CheckIpfsAlive(split[0], split[1]); !alive {
					model.DeleteFileIpfs(fileIpfs)
				}
			}
		}
	}
}

type SummaryResp struct {
	CidsCount  int64 `json:"cids_count"`
	Providers  int64 `json:"providers"`
	IpfsNodes  int64 `json:"ipfs_nodes"`
	DealsCount int64 `json:"deals_count"`
	DataStored int64 `json:"data_stored"`
	Height     int64 `json:"height"`
}

type FileSourcePager struct {
	Total     int64            `json:"total"`
	PageCount int64            `json:"pageCount"`
	Sources   []FileSourceResp `json:"sources"`
}

type FileSourceResp struct {
	FileName         string         `json:"file_name"`
	DataCid          string         `json:"data_cid"`
	FileSize         int64          `json:"file_size"`
	IpfsUrls         []string       `json:"ipfs_urls"`
	Providers        []ProviderInfo `json:"providers"`
	RebuildStatus    bool           `json:"rebuild_status"`
	HotBackups       int            `json:"hot_backups"`
	ColdBackups      int            `json:"cold_backups"`
	NotFoundProvider string         `json:"not_found_provider"`
	Status           string         `json:"status"`
}

type ProviderInfo struct {
	ProviderId string `json:"provider_id"`
	Status     string `json:"status"`
}

type IndexData struct {
	MultihashResults []struct {
		Multihash       string `json:"Multihash"`
		ProviderResults []struct {
			ContextID string `json:"ContextID"`
			Metadata  string `json:"Metadata"`
			Provider  struct {
				ID    string   `json:"ID"`
				Addrs []string `json:"Addrs"`
			} `json:"Provider"`
		} `json:"ProviderResults"`
	} `json:"MultihashResults"`
}

type RetrieveReq struct {
	DataCid string  `json:"data_cid"`
	CopyNum int     `json:"copy_num"`
	Cost    float64 `json:"cost"`
}
