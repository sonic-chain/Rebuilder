package model

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type OperateStatus string

const (
	START_PEER_STATUS       OperateStatus = "peerId-start"
	SUCCESS_PEER_STATUS     OperateStatus = "peerId-success"
	FAILED_PEER_STATUS      OperateStatus = "peerId-failed"
	START_MINER_STATUS      OperateStatus = "minerId-start"
	SUCCESS_MINER_STATUS    OperateStatus = "minerId-success"
	FAILED_MINER_STATUS     OperateStatus = "minerId-failed"
	START_RETRIEVE_STATUS   OperateStatus = "retrieve-start"
	SUCCESS_RETRIEVE_STATUS OperateStatus = "retrieve-success"
	FAILED_RETRIEVE_STATUS  OperateStatus = "retrieve-failed"
	START_IPFS_STATUS       OperateStatus = "ipfs-start"
	SUCCESS_IPFS_STATUS     OperateStatus = "ipfs-success"
	FAILED_IPFS_STATUS      OperateStatus = "ipfs-failed"
)

type MinerDeal struct {
	DealCid       string        `json:"deal_cid" gorm:"primaryKey"`
	DealId        int64         `json:"deal_id"`
	MinerId       string        `json:"miner_id"`
	Status        string        `json:"status"`
	PayloadCid    string        `json:"payload_cid"`
	PieceCid      string        `json:"piece_cid"`
	OperateStatus OperateStatus `json:"operate_status"`
}

func (MinerDeal) TableName() string {
	return "t_miner_deal"
}

func FindProvidersByPayloadCid(payloadCid string) ([]MinerDeal, error) {
	var minerDeals []MinerDeal
	if err := db.Where("payload_cid=?", payloadCid).Find(&minerDeals).Error; err != nil {
		return nil, err
	}
	return minerDeals, nil
}

func UpdateMinerDealStatus(payloadCid string, status OperateStatus) {
	if err := db.Model(&MinerDeal{}).Where("payload_cid=?", payloadCid).UpdateColumn("operate_status", status).Error; err != nil {
		log.Errorf("UpdateMinerDealStatus failed,payloadCid: %s, status: %s,error: %v", payloadCid, status, err)
	}
}

func DeleteMinerDeal(payloadCid string, minerIds []string) {
	if len(minerIds) > 0 {
		if err := db.Model(&MinerDeal{}).Where("miner_id not in (?)", minerIds).Delete(&MinerDeal{PayloadCid: payloadCid}).Error; err != nil {
			log.Errorf("DeleteMinerDeal failed,payloadCid: %s, minerIds: %v,error: %v", payloadCid, minerIds, err)
		}
	} else {
		if err := db.Debug().Model(&MinerDeal{}).Where("payload_cid=?", payloadCid).Delete(&MinerDeal{}).Error; err != nil {
			log.Errorf("DeleteMinerDeal failed,payloadCid: %s,error: %v", payloadCid, err)
		}
	}

}

type SourceFile struct {
	PayloadCid  string    `json:"payload_cid" gorm:"primaryKey"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	UploadId    int       `json:"upload_id"`
	Status      string    `json:"status"`
	Address     string    `json:"address"`
	CreateAt    time.Time `json:"create_at"`
	IndexStatus string    `json:"index_status"`
}

func (SourceFile) TableName() string {
	return "t_source_file"
}

func InsertSourceFile(sf *SourceFile) {
	db.Model(&SourceFile{}).Save([]*SourceFile{sf})
}

func FindSourceFileByStatus() ([]SourceFile, error) {
	var data []SourceFile
	result := db.Model(&SourceFile{}).Where("status in (?)", []string{"Processing", "Pending"}).Find(&data)
	if result.Error != nil {
		log.Errorf("find sourceFile by status=Processing failed, error: %v", result.Error)
		return data, result.Error
	}
	return data, nil
}

func FindSourceFile() ([]SourceFile, error) {
	var data []SourceFile
	result := db.Model(&SourceFile{}).Find(&data)
	if result.Error != nil {
		log.Errorf("find sourceFile by status=Processing failed, error: %v", result.Error)
		return data, result.Error
	}
	return data, nil
}

func UpdateSourceFileStatusByUploadId(dataCid string, uploadId int, status, indexStatus string) {
	result := db.Model(&SourceFile{}).Where("payload_cid=?", dataCid).Updates(map[string]interface{}{
		"status":       status,
		"index_status": indexStatus,
		"upload_id":    uploadId,
	})
	if result.Error != nil {
		log.Errorf("update sourceFile status failed,uploadId:%d, error: %v", uploadId, result.Error)
		return
	}
}

func FileSourceList(fieldName string, page int64, size int64) ([]SourceFile, error) {
	var fileList []SourceFile
	var err error
	if fieldName != "" {
		err = db.Model(&SourceFile{}).Where("payload_cid LIKE ?", "%"+fieldName+"%").Or("file_name LIKE ?", "%"+fieldName+"%").
			Order("create_at").Limit(int(size)).Offset(int(page * size)).Find(&fileList).Error
	} else {
		err = db.Model(&SourceFile{}).Order("create_at").Limit(int(size)).Offset(int(page * size)).Find(&fileList).Error
	}
	return fileList, err
}

func CountFileSourceList(fieldName string, size int64) (int64, int64, error) {
	var num, page int64
	var err error
	if fieldName != "" {
		err = db.Model(&SourceFile{}).Where("payload_cid LIKE ?", "%"+fieldName+"%").Or("file_name LIKE ?", "%"+fieldName+"%").
			Order("create_at").Count(&num).Error
	} else {
		err = db.Model(&SourceFile{}).Order("create_at").Count(&num).Error
	}
	if num%size == 0 {
		page = num / size
	} else {
		page = num/size + 1
	}
	return num, page, err
}

type FileIpfs struct {
	DataCid  string `json:"data_cid" gorm:"primaryKey"`
	IpfsUrl  string `json:"ipfs_url" gorm:"primaryKey"`
	IpfsHash string `json:"ipfs_hash"`
}

func (FileIpfs) TableName() string {
	return "t_file_ipfs"
}

func FindFileIpfsList() []FileIpfs {
	var result []FileIpfs
	db.Model(&FileIpfs{}).Where("ipfs_hash is not null").Find(&result)
	return result
}

func DeleteFileIpfs(fileIpfs FileIpfs) {
	db.Model(&FileIpfs{}).Delete(&fileIpfs)
}

func FindIpfsByDataCid(dataCid string) ([]string, error) {
	var ipfs []FileIpfs
	if err := db.Where("data_cid=?", dataCid).Find(&ipfs).Error; err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, fi := range ipfs {
		result = append(result, fi.IpfsUrl)
	}
	return result, nil
}

func FindIpfsCopysLow() ([]string, error) {
	result := make([]string, 0)
	// select payload_cid from t_source_file s left join t_file_ipfs f on s.payload_cid= f.data_cid group by s.payload_cid having sum(case when ipfs_url is null then 0 else 1 end) <1;
	if err := db.Model(&SourceFile{}).Select("payload_cid").
		Joins("left join t_file_ipfs on t_source_file.payload_cid=t_file_ipfs.data_cid").Where("t_source_file.index_status is null").
		Group("t_source_file.payload_cid").Having("sum(case when t_file_ipfs.ipfs_url is null then 0 else 1 end) < ?", 1).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func InsertFileIpfs(fileIpfs []FileIpfs) error {
	if err := db.Model(&FileIpfs{}).Save(fileIpfs).Error; err != nil {
		log.Errorf("insert FileIpfs data failed,error: %v", err)
		return err
	}
	return nil
}

type MinerPeer struct {
	MinerId string `json:"miner_id"`
	PeerId  string `json:"peer_id"`
}

func (MinerPeer) TableName() string {
	return "t_miner_peer"
}

func InsertMinerPeer(mp []MinerPeer) {
	if err := db.Model(&MinerPeer{}).CreateInBatches(mp, len(mp)).Error; err != nil {
		log.Errorf("insert minerpeer data failed,error: %v", err)
		return
	}
}

func DeleteMinerPeer() {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&MinerPeer{}).Error; err != nil {
		log.Errorf("delete minerpeer data failed,error: %v", err)
		return
	}
}

func FindMinerPeer(peerId string) (*MinerPeer, error) {
	var mp MinerPeer
	result := db.Model(&MinerPeer{}).First(&mp, MinerPeer{PeerId: peerId})
	if result.Error != nil {
		return &mp, result.Error
	}
	return &mp, nil
}

type FailedSource struct {
	PayloadCid string `json:"payload_cid" gorm:"primaryKey"`
}

func (FailedSource) TableName() string {
	return "t_failed_source"
}

func FindFailedSource() map[string]interface{} {
	var data []FailedSource
	result := make(map[string]interface{})
	db.Model(&FailedSource{}).Find(&data)
	for _, source := range data {
		result[source.PayloadCid] = struct{}{}
	}
	return result
}

func InsertFailedSource(fs FailedSource) {
	db.Model(&FailedSource{}).Create(&fs)
}

func DeleteFailedSource(fs FailedSource) {
	db.Model(&FailedSource{}).Delete(&fs)
}

func CountFileSource() (int64, error) {
	var num int64
	var err error
	err = db.Model(&SourceFile{}).Count(&num).Error
	return num, err
}

func CountDealByMinerDeal() (int64, error) {
	var num int64
	var err error
	err = db.Model(&MinerDeal{}).Count(&num).Error
	return num, err
}

func CountProviderMinerDeal() (int64, error) {
	var num int64
	var err error
	err = db.Model(&MinerDeal{}).Group("miner_id").Count(&num).Error
	return num, err
}

func SaveOrUpdateMinerDeal(md []MinerDeal) error {
	return db.Model(MinerDeal{}).Save(md).Error
}

func IpfsNodeCount() (int64, error) {
	var fIpfs []FileIpfs
	var err error
	if err = db.Model(&FileIpfs{}).Find(&fIpfs).Error; err != nil {
		return 0, err
	}

	urlsMap := make(map[string]struct{})
	for _, f := range fIpfs {
		splits := strings.Split(f.IpfsUrl, "/")
		urlsMap[splits[2]] = struct{}{}
	}
	return int64(len(urlsMap)), nil
}

func HotDataSize() ([]DataSizeResult, error) {
	result := make([]DataSizeResult, 0)
	// select t_source_file.payload_cid,file_size,sum(case when ipfs_url is null then 0 else 1 end) as ipfsSum from t_source_file  left join t_file_ipfs  on t_source_file.payload_cid= t_file_ipfs.data_cid group by t_source_file.payload_cid;
	if err := db.Model(&SourceFile{}).Select("t_source_file.payload_cid,t_source_file.file_size,sum(case when ipfs_url is null then 0 else 1 end) as num").
		Joins("left join t_file_ipfs on t_source_file.payload_cid= t_file_ipfs.data_cid").
		Group("t_source_file.payload_cid").Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func ColdDataSize() ([]DataSizeResult, error) {
	result := make([]DataSizeResult, 0)
	// select t_source_file.payload_cid,file_size,count(miner_id) as num from t_source_file  left join t_miner_deal  on t_source_file.payload_cid= t_miner_deal.payload_cid group by t_source_file.payload_cid;
	if err := db.Model(&SourceFile{}).Select("t_source_file.payload_cid,t_source_file.file_size,count(miner_id) as num").
		Joins("left join t_miner_deal on t_source_file.payload_cid=t_miner_deal.payload_cid").
		Group("t_source_file.payload_cid").Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

type DataSizeResult struct {
	PayloadCid string `json:"payload_cid"`
	FileSize   int64  `json:"file_size"`
	Num        int64  `json:"num"`
}
