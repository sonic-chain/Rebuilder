package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/common"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/internal"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/model"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
	"github.com/unknwon/com"
)

var log = logging.Logger("service")

// @Summary 汇总信息展示
// @Produce  json
// @Success 200 {object} internal.Response{data=service.SummaryResp} "正常返回的数据格式"
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
		log.Errorf("Summary get IpfsNodeCount failed,error: %v", err)
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

// @Summary 更新rebuilder任务的状态
// @Accept  json
// @Param   data	body	service.RebuildStatusReq  true	"请求参数"
// @Produce  json
// @Success 200 {object} internal.Response
// @Failure 500 {object} internal.Response
// @Router /rebuild/status [post]
func RebuildStatus(c *gin.Context) {
	appG := internal.Gin{C: c}
	var rebuildStatus RebuildStatusReq
	if err := c.ShouldBindJSON(&rebuildStatus); err != nil {
		appG.Response(http.StatusBadRequest, internal.INVALID_PARAMS, internal.GetMsg(internal.INVALID_PARAMS))
	}

	model.InsertFileIpfs([]model.FileIpfs{{
		DataCid: rebuildStatus.DataCid,
		IpfsUrl: rebuildStatus.IpfsUrl,
	}})

	model.UpdateSourceFileStatusByUploadId(rebuildStatus.DataCid, rebuildStatus.SourceFileUploadId, "Processing", "")

	appG.Response(http.StatusOK, internal.SUCCESS, "")
}

// @Summary 获取文件存储信息列表
// @Produce  json
// @param	field_name	query	string	false	"data_cid/file_name"
// @param	page		query	int		false	"页码，默认从0开始"
// @param	size		query	int		false	"条数，默认为20条"
// @Success 200 {object} internal.Response{data=service.FileSourcePager} "正常返回的数据格式"
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

	failedSources := model.FindFailedSource()
	sourceList, err := model.FileSourceList(fieldName, page, size)
	if err != nil {
		appG.Response(http.StatusInternalServerError, internal.ERROR_FILE_LIST_FAIL, nil)
		return
	}

	result := make([]FileSourceResp, 0)
	for _, fileSource := range sourceList {
		var fsr FileSourceResp
		fsr.FileName = fileSource.FileName
		fsr.DataCid = fileSource.PayloadCid
		fsr.FileSize = fileSource.FileSize
		fsr.McsStatus = fileSource.Status
		fsr.UploadId = fileSource.UploadId

		ipfsUrls, err := model.FindIpfsByDataCid(fileSource.PayloadCid)
		if err != nil {
			log.Errorf("FindIpfsBySourceId failed, datacid: %s,error: %v", fileSource.PayloadCid, err)
			continue
		}
		fsr.IpfsUrls = ipfsUrls
		fsr.HotBackups = len(ipfsUrls)

		if _, ok := failedSources[fileSource.PayloadCid]; !ok {
			providers, err := model.FindProvidersByPayloadCid(fileSource.PayloadCid)
			if err != nil {
				log.Errorf("FindProvidersByPayloadCid failed, payload_cid: %s,error: %v", fileSource.PayloadCid, err)
				continue
			}
			providerStatus := make([]ProviderInfo, 0)
			for _, p := range providers {
				providerStatus = append(providerStatus, ProviderInfo{
					ProviderId: p.MinerId,
					Status:     p.Status,
				})
			}
			fsr.Providers = providerStatus
			fsr.ColdBackups = len(providers)
		} else {
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

// @Summary 根据cid从filecoin检索文件存储到ipfs
// @Produce  json
// @Param cid path string true "ID"
// @Success 200 {object} internal.Response
// @Failure 500 {object} internal.Response
// @Router /cid/{cid} [get]
func GetCid(c *gin.Context) {
	appG := internal.Gin{C: c}
	cid := com.StrTo(c.Param("cid")).String()
	//  1. query peerId by indexer node
	model.UpdateMinerDealStatus(cid, model.START_PEER_STATUS)
	peerData := common.NewIndexerClient().SendHttpGet(common.GET_PEER_URL, cid)

	peerIds := make(map[string]string, 0)
	for _, data := range peerData {
		if string(data) == "no results for query" {
			continue
		}
		var indexData IndexData
		if err := json.Unmarshal(data, &indexData); err != nil {
			log.Errorf("change to json failed,error: %v", err)
			model.UpdateMinerDealStatus(cid, model.FAILED_PEER_STATUS)
			continue
		}
		if len(indexData.MultihashResults) > 0 && len(indexData.MultihashResults[0].ProviderResults) > 0 {
			peerId := indexData.MultihashResults[0].ProviderResults[0].Provider.ID
			peerIds[peerId] = peerId
		}
	}
	if len(peerIds) == 0 {
		model.UpdateMinerDealStatus(cid, model.FAILED_PEER_STATUS)
		model.InsertFailedSource(model.FailedSource{
			PayloadCid: cid,
		})
		appG.Response(http.StatusInternalServerError, internal.ERROR_RETRIEVE_FAIL, nil)
		return
	} else {
		model.UpdateMinerDealStatus(cid, model.SUCCESS_PEER_STATUS)
		model.DeleteFailedSource(model.FailedSource{
			PayloadCid: cid,
		})
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("catch panic error message: %v \n", err)
			}
		}()

		lotusClient := common.NewLotusClient()
		for _, peerId := range peerIds {
			log.Infof("start process peerId: %s", peerId)
			// 2. query minerId by lotus
			model.UpdateMinerDealStatus(cid, model.START_MINER_STATUS)

			minerId, err := lotusClient.GetMinerIdByPeerId(peerId)
			if err != nil {
				log.Warnf("get minerpeer failed,peerId:%s,error: %v,continue check next peerId", peerId, err)
				model.UpdateMinerDealStatus(cid, model.FAILED_MINER_STATUS)
				continue
			}
			model.UpdateMinerDealStatus(cid, model.SUCCESS_MINER_STATUS)

			savePath := filepath.Join(model.LotusSetting.DownloadDir, minerId+"-"+cid)
			// 3. retrieveData
			model.UpdateMinerDealStatus(cid, model.START_RETRIEVE_STATUS)
			if err := lotusClient.RetrieveData(minerId, cid, savePath); err != nil {
				model.UpdateMinerDealStatus(cid, model.FAILED_RETRIEVE_STATUS)
				continue
			}
			model.UpdateMinerDealStatus(cid, model.SUCCESS_RETRIEVE_STATUS)

			// 4. upload file to ipfs  ERROR_UPLOAD_FAIL
			model.UpdateMinerDealStatus(cid, model.START_IPFS_STATUS)

			stat, err := os.Stat(savePath)
			if err != nil {
				log.Errorf("not found savepath: %s,error: %s", savePath, err)
				return
			}

			fileIpfs := make([]model.FileIpfs, 0)
			if stat.IsDir() {
				urls := UploaderDir(savePath, model.UploaderSetting.IpfsUrls)
				for _, u := range urls {
					split := strings.Split(u, "/")
					fileIpfs = append(fileIpfs, model.FileIpfs{
						DataCid:  cid,
						IpfsUrl:  u,
						IpfsHash: split[len(split)-1],
					})
				}
			} else {
				hashs, urls := UploaderFile(savePath, model.UploaderSetting.IpfsUrls)
				for index, u := range urls {
					fileIpfs = append(fileIpfs, model.FileIpfs{
						DataCid:  cid,
						IpfsUrl:  u,
						IpfsHash: hashs[index],
					})
				}
			}
			if len(fileIpfs) > 0 {
				if err = model.InsertFileIpfs(fileIpfs); err == nil {
					model.UpdateMinerDealStatus(cid, model.SUCCESS_IPFS_STATUS)
					err := os.RemoveAll(savePath)
					if err != nil {
						log.Errorf("remove file failed savepath: %s,error: %v", savePath, err)
					}
				}
			}
		}
	}()
	appG.Response(http.StatusOK, internal.SUCCESS, map[string]interface{}{
		"msg": "已经提交处理中",
	})
}

func AutoUploadFileToIpfs() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			payloadCids, err := model.FindIpfsCopysLow()
			if err != nil {
				return
			}
			for _, sid := range payloadCids {
				cid := sid
				model.UpdateMinerDealStatus(cid, model.START_PEER_STATUS)
				peerData := common.NewIndexerClient().SendHttpGet(common.GET_PEER_URL, cid)

				peerIds := make(map[string]string, 0)
				for _, data := range peerData {
					if string(data) == "no results for query" {
						continue
					}
					var indexData IndexData
					if err = json.Unmarshal(data, &indexData); err != nil {
						model.UpdateMinerDealStatus(cid, model.FAILED_PEER_STATUS)
						continue
					}
					if len(indexData.MultihashResults) > 0 && len(indexData.MultihashResults[0].ProviderResults) > 0 {
						peerId := indexData.MultihashResults[0].ProviderResults[0].Provider.ID
						peerIds[peerId] = peerId
					}
				}
				if len(peerIds) == 0 {
					model.UpdateMinerDealStatus(cid, model.FAILED_PEER_STATUS)
					model.InsertFailedSource(model.FailedSource{
						PayloadCid: cid,
					})
					model.UpdateSourceFileStatusByUploadId(cid, 0, "", "failed")
					continue
				} else {
					model.UpdateMinerDealStatus(cid, model.SUCCESS_PEER_STATUS)
					model.DeleteFailedSource(model.FailedSource{
						PayloadCid: cid,
					})
				}

				go func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Printf("catch panic error message: %v \n", err)
						}
					}()

					lotusClient := common.NewLotusClient()
					for _, peerId := range peerIds {
						log.Infof("start process peerId: %s", peerId)
						// 2. query minerId by lotus
						model.UpdateMinerDealStatus(cid, model.START_MINER_STATUS)

						minerId, err := lotusClient.GetMinerIdByPeerId(peerId)
						if err != nil {
							log.Warnf("get minerpeer failed,peerId:%s,error: %v,continue check next peerId", peerId, err)
							model.UpdateMinerDealStatus(cid, model.FAILED_MINER_STATUS)
							continue
						}
						model.UpdateMinerDealStatus(cid, model.SUCCESS_MINER_STATUS)

						savePath := filepath.Join(model.LotusSetting.DownloadDir, minerId+"-"+cid)
						// 3. retrieveData
						model.UpdateMinerDealStatus(cid, model.START_RETRIEVE_STATUS)
						if err := lotusClient.RetrieveData(minerId, cid, savePath); err != nil {
							model.UpdateMinerDealStatus(cid, model.FAILED_RETRIEVE_STATUS)
							continue
						}
						model.UpdateMinerDealStatus(cid, model.SUCCESS_RETRIEVE_STATUS)

						// 4. upload file to ipfs  ERROR_UPLOAD_FAIL
						model.UpdateMinerDealStatus(cid, model.START_IPFS_STATUS)

						stat, err := os.Stat(savePath)
						if err != nil {
							log.Errorf("not found savepath: %s,error: %s", savePath, err)
							return
						}

						fileIpfs := make([]model.FileIpfs, 0)
						if stat.IsDir() {
							urls := UploaderDir(savePath, model.UploaderSetting.IpfsUrls)
							for _, u := range urls {
								split := strings.Split(u, "/")
								fileIpfs = append(fileIpfs, model.FileIpfs{
									DataCid:  cid,
									IpfsUrl:  u,
									IpfsHash: split[len(split)-1],
								})
							}
						} else {
							hashs, urls := UploaderFile(savePath, model.UploaderSetting.IpfsUrls)
							for index, u := range urls {
								fileIpfs = append(fileIpfs, model.FileIpfs{
									DataCid:  cid,
									IpfsUrl:  u,
									IpfsHash: hashs[index],
								})
							}
						}
						if len(fileIpfs) > 0 {
							if err = model.InsertFileIpfs(fileIpfs); err == nil {
								model.UpdateMinerDealStatus(cid, model.SUCCESS_IPFS_STATUS)
								err := os.RemoveAll(savePath)
								if err != nil {
									log.Errorf("remove file failed savepath: %s,error: %v", savePath, err)
									model.UpdateMinerDealStatus(cid, model.FAILED_IPFS_STATUS)
								}
							}
						}

					}
				}()
			}
		}
	}
}

func AutoSourceFileStatusAndMinerDealInfo() {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			sourceFiles, err := model.FindSourceFileByStatus()
			if err != nil {
				return
			}
			for _, sf := range sourceFiles {
				mcsStatus, err := common.GetMcsStatus(sf.FileName, sf.Address)
				if err != nil {
					log.Errorf("get mcs status failed,error: %v", err)
					continue
				}
				if mcsStatus.Status == "success" {
					for _, mcs := range mcsStatus.Data.SourceFileUpload {
						if mcs.SourceFileUploadId == sf.UploadId {
							if mcs.DealSuccess {
								model.UpdateSourceFileStatusByUploadId(sf.PayloadCid, mcs.SourceFileUploadId, "Success", "")
								model.DeleteFailedSource(model.FailedSource{
									PayloadCid: sf.PayloadCid,
								})
							} else {
								model.UpdateSourceFileStatusByUploadId(sf.PayloadCid, mcs.SourceFileUploadId, mcs.Status, "")
							}
							minerDeals := make([]model.MinerDeal, 0)
							for _, deal := range mcs.OfflineDeal {
								var md model.MinerDeal
								md.DealCid = deal.DealCid
								md.DealId = int64(deal.DealId)
								md.MinerId = deal.MinerFid
								md.Status = deal.OnChainStatus
								md.PayloadCid = sf.PayloadCid
								minerDeals = append(minerDeals, md)
							}
							if len(minerDeals) > 0 {
								model.SaveOrUpdateMinerDeal(minerDeals)
							}
						}
					}
				}
			}
		}
	}
}

func WatchFilecoinNodeData() {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			sourceFiles, err := model.FindSourceFile()
			if err != nil {
				return
			}
			for _, sourceFile := range sourceFiles {
				sourceFileCp := sourceFile
				peerData := common.NewIndexerClient().SendHttpGet(common.GET_PEER_URL, sourceFileCp.PayloadCid)

				peerIds := make(map[string]string, 0)
				for _, data := range peerData {
					if string(data) == "no results for query" {
						continue
					}
					var indexData IndexData
					if err = json.Unmarshal(data, &indexData); err != nil {
						continue
					}
					if len(indexData.MultihashResults) > 0 && len(indexData.MultihashResults[0].ProviderResults) > 0 {
						peerId := indexData.MultihashResults[0].ProviderResults[0].Provider.ID
						peerIds[peerId] = peerId
					}
				}

				updateMinerIds := make([]string, 0)
				lotusClient := common.NewLotusClient()
				for _, peerId := range peerIds {
					minerId, err := lotusClient.GetMinerIdByPeerId(peerId)
					if err != nil {
						log.Warnf("get minerpeer failed,peerId:%s,error: %v,continue check next peerId", peerId, err)
						continue
					}
					updateMinerIds = append(updateMinerIds, minerId)
				}
				model.DeleteMinerDeal(sourceFileCp.PayloadCid, updateMinerIds)
			}
		}
	}
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
	HotBackups       int            `json:"hot_backups"`
	McsStatus        string         `json:"mcs_status"`
	ColdBackups      int            `json:"cold_backups"`
	UploadId         int            `json:"upload_id"`
	NotFoundProvider string         `json:"not_found_provider"`
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

type RebuildStatusReq struct {
	SourceFileUploadId int    `json:"source_file_upload_id"`
	PayloadCid         string `json:"payload_cid"`
	IpfsUrl            string `json:"ipfs_url"`
	FileSize           int    `json:"file_size"`
	DataCid            string `json:"data_cid"`

	//FileName           string `json:"file_name"`
	//PinStatus          string `json:"pin_status"`
	//Status             string `json:"status"`
	//OfflineDeal        []Deal `json:"offline_deal"`
}

type Deal struct {
	DealCid       string `json:"deal_cid"`
	Status        string `json:"status"`
	DealId        int64  `json:"deal_id"`
	OnChainStatus string `json:"on_chain_status"`
	MinerFid      string `json:"miner_fid"`
}
