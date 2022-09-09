package model

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var log = logging.Logger("model")

func SetupDB() {
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", DatabaseSetting.User,
			DatabaseSetting.Password, DatabaseSetting.Host, DatabaseSetting.Name), // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		CreateBatchSize: 100,
	})

	if err != nil {
		log.Fatalf("models.SetupDB err: %v", err)
	}

	db.AutoMigrate(
		&MinerDeal{},
		&SourceFile{},
		&FileIpfs{},
		&MinerPeer{},
	)
}
