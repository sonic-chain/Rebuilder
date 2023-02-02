package main

import (
	"fmt"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/common"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/model"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/routers"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/service"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
	"syscall"
)

var log = logging.Logger("main")

// @title ReBuilder API
// @version 1.0
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @host      api.storefrontiers.cn
// @BasePath  /api/v1
func main() {
	lvl, err := logging.LevelFromString("info")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)
	common.InitConfig()
	model.SetupDB()
	gin.SetMode(model.ServerSetting.RunMode)
	routersInit := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", model.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20
	go service.AutoUploadFileToIpfs()
	go service.AutoSourceFileStatusAndMinerDealInfo()
	go service.WatchIpfsNodeData()
	//go service.WatchFilecoinNodeData()
	endless.DefaultMaxHeaderBytes = maxHeaderBytes
	server := endless.NewServer(endPoint, routersInit)
	server.BeforeBegin = func(add string) {
		log.Infof("Actual pid is %d", syscall.Getpid())
	}
	log.Infof("[info] start http server listening %s", endPoint)
	if err = server.ListenAndServe(); err != nil {
		log.Infof("Server err: %v", err)
	}
}

//func main() {
//
//	localFileName := "/Users/sonic/Documents/go_work/Rebuilder/李信-王者荣耀.jpg"
//	filePaths := strings.Split(localFileName, "/")
//	objectName := path.Join(time.Now().Format("2006-01-02"), filePaths[len(filePaths)-1])
//
//	uploadFile, err := mcs.UploadFile(context.TODO(), "rebuilder", objectName, localFileName)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	fmt.Printf("%+v \n", string(uploadFile))
//
//	ipfsUrl, err := mcs.GetFile(context.TODO(), "rebuilder", objectName)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	fmt.Printf("ipfsUrl: %s \n", ipfsUrl)
//
//	buckets, err := mcs.GetBuckets(context.TODO())
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	fmt.Printf("buckets: %+v \n", buckets)
//
//}
