package common

import (
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/model"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
	model.ContractConfig.Address = viper.GetString("contract.address")
	model.ContractConfig.RpcUrl = viper.GetString("contract.rpcUrl")

}
