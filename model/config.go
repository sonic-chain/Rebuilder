package model

type Server struct {
	RunMode  string
	HttpPort int
}

var ServerSetting = &Server{}

type Database struct {
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Indexer struct {
	Urls []string
}

var IndexerSetting = &Indexer{}

type lotus struct {
	FullNodeApi string
	DownloadDir string
	Address     string
}

var LotusSetting = &lotus{}

type uploader struct {
	IpfsUrls []string
}

var UploaderSetting = &uploader{}

type Contract struct {
	Address string
	RpcUrl  string
	Private string
}

var ContractConfig = &Contract{}
