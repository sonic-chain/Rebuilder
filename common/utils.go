package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/model"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: [%s]", viper.ConfigFileUsed())
	} else {
		log.Errorf("read config file failed,error: %v", err)
	}
	setConfig()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		setConfig()
		log.Infof("Config file changed: [%s]", e.Name)
	})
}

func setConfig() {
	model.ServerSetting.RunMode = viper.GetString("server.RunMode")
	model.ServerSetting.HttpPort = viper.GetInt("server.HttpPort")

	model.DatabaseSetting.User = viper.GetString("database.User")
	model.DatabaseSetting.Password = viper.GetString("database.Password")
	model.DatabaseSetting.Host = viper.GetString("database.Host")
	model.DatabaseSetting.Name = viper.GetString("database.Name")

	model.IndexerSetting.Urls = viper.GetStringSlice("indexer.Urls")
	model.LotusSetting.FullNodeApi = viper.GetString("lotus.FullNodeApi")
	model.LotusSetting.DownloadDir = viper.GetString("lotus.DownloadDir")
	model.LotusSetting.Address = viper.GetString("lotus.Address")
	model.UploaderSetting.IpfsUrls = viper.GetStringSlice("uploader.IpfsUrls")
}

func SaveMinerId() error {
	minerFile, err := os.Open("./miner.csv")
	if err != nil {
		return err
	}
	defer minerFile.Close()

	reader := bufio.NewReader(minerFile)
	isFirst := true
	mp := make([]model.MinerPeer, 0)
	var count int
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			if isFirst {
				isFirst = false
				continue
			}
			count++
			line = strings.Trim(line, "\n")
			data := strings.Split(line, ",")
			fmt.Printf("minerId: %s, peerID:[%s] \n", data[0], data[1])
			if data[0] == "" || data[1] == "" {
				continue
			}
			mp = append(mp, model.MinerPeer{
				MinerId: data[0],
				PeerId:  data[1],
			})
			if count > 50 {
				model.InsertMinerPeer(mp)
				mp = make([]model.MinerPeer, 0)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("read miner.csv, error: %+v", err)
		}
		if len(mp) > 0 {
			model.InsertMinerPeer(mp)
		}
	}
	return nil
}

func CreateTableData() {
	for i := 0; i < 100; i++ {
		var url = fmt.Sprintf("https://mcs-api.filswan.com/api/v1/storage/tasks/deals?page_size=10&page_number=%d&wallet_address=0x3350bfbcd9ac435cd3c410bc98e1ec5b94a662e5", i)
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("get mcs status failed, url: %s,error: %s", url, err)

		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("read mcs status  resp failed,error: %s", err)
		}
		var mcs McsStatus
		if err := json.Unmarshal(data, &mcs); err != nil {
		}

		if mcs.Status == "success" && len(mcs.Data.SourceFileUpload) == 0 {
			break
		}
		if mcs.Status == "success" {
			for _, mcs := range mcs.Data.SourceFileUpload {
				var sf model.SourceFile
				sf.PayloadCid = mcs.PayloadCid
				sf.FileName = mcs.FileName
				sf.FileSize = int64(mcs.FileSize)
				sf.Status = mcs.Status
				rand.Seed(time.Now().UnixNano())
				num := rand.Intn(30000)
				sf.CreateAt = time.Now().Add(time.Duration(-num) * time.Second)
				model.InsertSourceFile(&sf)

				model.InsertFileIpfs([]model.FileIpfs{{
					DataCid: mcs.PayloadCid,
					IpfsUrl: mcs.IpfsUrl,
				}})

				minerDeals := make([]model.MinerDeal, 0)
				for _, deal := range mcs.OfflineDeal {
					var md model.MinerDeal
					md.DealCid = deal.DealCid
					md.DealId = int64(deal.DealId)
					md.MinerId = deal.MinerFid
					md.Status = deal.OnChainStatus
					md.PayloadCid = mcs.PayloadCid
					minerDeals = append(minerDeals, md)
				}
				if len(minerDeals) > 0 {
					model.SaveOrUpdateMinerDeal(minerDeals)
				}
			}
		}
	}
}

func GetMcsStatus(fileName, wallet string) (McsStatus, error) {
	var url = fmt.Sprintf("https://mcs-api.filswan.com/api/v1/storage/tasks/deals?page_size=20&page_number=1&wallet_address=%s&file_name=%s", wallet, fileName)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("get mcs status failed, url: %s,error: %s", url, err)
		return McsStatus{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read mcs status  resp failed,error: %s", err)
		return McsStatus{}, err
	}
	var mcs McsStatus
	if err := json.Unmarshal(data, &mcs); err != nil {
		return McsStatus{}, err
	}
	return mcs, nil
}

type McsStatus struct {
	Status string `json:"status"`
	Data   struct {
		SourceFileUpload []SourceFileUpload `json:"source_file_upload"`
		TotalRecordCount int                `json:"total_record_count"`
	} `json:"data"`
}

type SourceFileUpload struct {
	SourceFileUploadId int               `json:"source_file_upload_id"`
	CarFileId          int               `json:"car_file_id"`
	FileName           string            `json:"file_name"`
	FileSize           int               `json:"file_size"`
	UploadAt           int               `json:"upload_at"`
	Duration           int               `json:"duration"`
	IpfsUrl            string            `json:"ipfs_url"`
	PinStatus          string            `json:"pin_status"`
	PayloadCid         string            `json:"payload_cid"`
	WCid               string            `json:"w_cid"`
	Status             string            `json:"status"`
	DealSuccess        bool              `json:"deal_success"`
	OfflineDeal        []OfflineDealInfo `json:"offline_deal"`
}

type OfflineDealInfo struct {
	Id             int    `json:"id"`
	CarFileId      int    `json:"car_file_id"`
	DealCid        string `json:"deal_cid"`
	MinerId        int    `json:"miner_id"`
	Verified       bool   `json:"verified"`
	StartEpoch     int    `json:"start_epoch"`
	SenderWalletId int    `json:"sender_wallet_id"`
	Status         string `json:"status"`
	DealId         int    `json:"deal_id"`
	OnChainStatus  string `json:"on_chain_status"`
	UnlockTxHash   string `json:"unlock_tx_hash"`
	UnlockAt       int    `json:"unlock_at"`
	CreateAt       int    `json:"create_at"`
	UpdateAt       int    `json:"update_at"`
	MinerFid       string `json:"miner_fid"`
}
