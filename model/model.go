package model

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type RebuildStatus string

const (
	REBUILD_INDEXING          RebuildStatus = "Indexing"
	REBUILD_INDEXING_FAILED   RebuildStatus = "Indexing-failed"
	REBUILD_RETRIEVING        RebuildStatus = "Retrieving"
	REBUILD_RETRIEVING_FAILED RebuildStatus = "Retrieving-failed"
	REBUILD_UPLOADING         RebuildStatus = "Uploading"
	REBUILD_UPLOADING_FAILED  RebuildStatus = "Uploading-failed"
	REBUILD_SUCCESS           RebuildStatus = "Successful"
	REBUILD_SUCCESS_FAILED    RebuildStatus = "Successful-failed"

	BUCKET_NAME string = "rebuilder"
)

type SourceFile struct {
	DataCid       string        `json:"data_cid" gorm:"primaryKey"`
	FileName      string        `json:"file_name"`
	FileSize      int64         `json:"file_size"`
	RebuildFlag   bool          `json:"rebuild_flag"`
	RebuildStatus RebuildStatus `json:"rebuild_status"`
	IpfsUrls      []FileIpfs    `json:"ipfs_urls" gorm:"foreignKey:DataCid"`
	MinerIds      []FileMiner   `json:"miner_ids" gorm:"foreignKey:DataCid"`
	CreateAt      time.Time     `json:"create_at"`
}

func (SourceFile) TableName() string {
	return "t_source_file"
}

func InsertSourceFile(sf *SourceFile) {
	db.Model(&SourceFile{}).Save([]*SourceFile{sf})
}

type FileIpfs struct {
	DataCid string `json:"data_cid" gorm:"primaryKey"`
	IpfsUrl string `json:"ipfs_url"`
}

func (FileIpfs) TableName() string {
	return "t_file_ipfs"
}

type FileMiner struct {
	gorm.Model
	DataCid string `json:"data_cid"`
	MinerId string `json:"miner_id"`
	Status  string `json:"status"`
}

func (FileMiner) TableName() string {
	return "t_file_miner"
}

func UpdateSourceFileStatus(dataCid string, status RebuildStatus) {
	result := db.Model(&SourceFile{}).Where("data_cid = ?", dataCid).Updates(SourceFile{
		RebuildStatus: status,
	})
	if result.Error != nil {
		log.Errorf("update sourceFile status failed, data_cid: %s, error: %v", dataCid, result.Error)
		return
	}
}

func FileSourceList(fieldName string, page int64, size int64) ([]SourceFile, error) {
	var fileList []SourceFile
	var err error
	if fieldName != "" {
		err = db.Model(&SourceFile{}).Where("data_cid LIKE ?", "%"+fieldName+"%").Or("file_name LIKE ?", "%"+fieldName+"%").
			Order("create_at").Preload("IpfsUrls").Preload("MinerIds").Limit(int(size)).Offset(int(page * size)).Find(&fileList).Error
	} else {
		err = db.Model(&SourceFile{}).Order("create_at").Preload("IpfsUrls").Preload("MinerIds").Limit(int(size)).Offset(int(page * size)).Find(&fileList).Error
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

func InsertFileIpfs(fileIpfs []FileIpfs) error {
	if err := db.Model(&FileIpfs{}).Save(fileIpfs).Error; err != nil {
		log.Errorf("insert FileIpfs data failed,error: %v", err)
		return err
	}
	return nil
}

type Miner struct {
	MinerId string `json:"miner_id" gorm:"primaryKey"`
}

func (Miner) TableName() string {
	return "t_miner"
}

func InsertMiners(miners []Miner) {
	if err := db.Model(&MinerPeer{}).CreateInBatches(miners, len(miners)).Error; err != nil {
		log.Errorf("insert minerpeer data failed,error: %v", err)
		return
	}
}

func GetMiners() []Miner {
	var miners []Miner
	if err := db.Model(&Miner{}).Find(&miners).Error; err != nil {
		log.Errorf("get miner data failed,error: %v", err)
	}
	return miners
}

type MinerPeer struct {
	MinerId string `json:"miner_id"`
	PeerId  string `json:"peer_id"`
}

func (MinerPeer) TableName() string {
	return "t_miner_peer"
}

func InsertMinerPeers(mp []MinerPeer) {
	if err := db.Model(&MinerPeer{}).CreateInBatches(mp, len(mp)).Error; err != nil {
		log.Errorf("insert minerpeer data failed,error: %v", err)
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

func CountFileSource() (int64, error) {
	var num int64
	var err error
	err = db.Model(&SourceFile{}).Count(&num).Error
	return num, err
}

func CountDealByMinerDeal() (int64, error) {
	var num int64
	var err error
	err = db.Model(&FileMiner{}).Count(&num).Error
	return num, err
}

func CountProviderMinerDeal() (int64, error) {
	var num int64
	var err error
	err = db.Model(&FileMiner{}).Group("miner_id").Count(&num).Error
	return num, err
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
	// select t_source_file.data_cid,file_size,sum(case when ipfs_url is null then 0 else 1 end) as ipfsSum from t_source_file left join t_file_ipfs on t_source_file.data_cid=t_file_ipfs.data_cid group by t_source_file.data_cid;
	if err := db.Model(&SourceFile{}).Select("t_source_file.data_cid,t_source_file.file_size,sum(case when ipfs_url is null then 0 else 1 end) as num").
		Joins("left join t_file_ipfs on t_source_file.data_cid= t_file_ipfs.data_cid").
		Group("t_source_file.data_cid").Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func ColdDataSize() ([]DataSizeResult, error) {
	result := make([]DataSizeResult, 0)
	// select t_source_file.data_cid,file_size,count(miner_id) as num from t_source_file left join t_file_miner on t_source_file.data_cid=t_file_miner.data_cid group by t_source_file.data_cid;
	if err := db.Model(&SourceFile{}).Select("t_source_file.data_cid,t_source_file.file_size,count(miner_id) as num").
		Joins("left join t_file_miner on t_source_file.data_cid=t_file_miner.data_cid").
		Group("t_source_file.data_cid").Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func FindFileIpfsList() []FileIpfs {
	var result []FileIpfs
	db.Model(&FileIpfs{}).Find(&result)
	return result
}

func DeleteFileIpfs(fileIpfs FileIpfs) {
	db.Model(&FileIpfs{}).Delete(&fileIpfs)
}

type DataSizeResult struct {
	PayloadCid string `json:"payload_cid"`
	FileSize   int64  `json:"file_size"`
	Num        int64  `json:"num"`
}
