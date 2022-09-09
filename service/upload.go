package service

import (
	"strings"

	"github.com/filswan/go-swan-lib/client/ipfs"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"

	goipfs "github.com/ipfs/go-ipfs-api"
)

// targetFile 包含全路径的目标文件名
// h 上传成功后，ipfs节点返回的文件cid
// u 上传成功后，ipfs节点下载链接
func UploaderFile(targetFile string, ipfsApiUrls []string) (h []string, u []string) {
	//fmt.Println("Uploader Start...")

	hash := make([]string, 0)
	dlUrl := make([]string, 0)
	for index := 0; index < len(ipfsApiUrls); index++ {
		u := strings.Split(ipfsApiUrls[index], ";")
		uploadUrl := u[0]
		downloadUrl := u[1]

		apiUrlFull := utils.UrlJoin(uploadUrl, "api/v0/add?stream-channels=true&pin=true")
		ipfsFileHash, err := ipfs.IpfsUploadFileByStream(apiUrlFull, targetFile)
		if err != nil {
			logs.GetLogger().Error("Uploader IpfsUploadFileByStream: apiUrlFull=", apiUrlFull, "err=", err)
			continue
		}

		ipfsDownloadUrl := utils.UrlJoin(downloadUrl, "/ipfs/"+ipfsFileHash)
		logs.GetLogger().Info("Uploader Result:", "TargetFile=", targetFile, " Hash=", ipfsFileHash, "DownloadUrl=", ipfsDownloadUrl)

		hash = append(hash, ipfsFileHash)
		dlUrl = append(dlUrl, ipfsDownloadUrl)
	}

	return hash, dlUrl
}

func UploaderDir(targetDir string, ipfsApiUrls []string) (h []string) {

	dirHash := make([]string, 0)
	for index := 0; index < len(ipfsApiUrls); index++ {
		u := strings.Split(ipfsApiUrls[index], ";")
		uploadUrl := u[0]
		downloadUrl := u[1]

		cl := goipfs.NewShell(uploadUrl)
		hash, err := cl.AddDir(targetDir)
		if err != nil {
			logs.GetLogger().Error("UploaderDir:", " uploadUrl=", uploadUrl, " targetDir=", targetDir, " err=", err)
			continue
		}

		ipfsDownloadUrl := utils.UrlJoin(downloadUrl, "/ipfs/"+hash)
		logs.GetLogger().Info("Uploader Result:", "TargetDir=", targetDir, " Hash=", hash, "DownloadUrl=", ipfsDownloadUrl)

		dirHash = append(dirHash, ipfsDownloadUrl)
	}

	return dirHash
}

//downUrl  下载url
//tarFile  下载文件保存 路径/文件名
func Downloader(downUrl string, tarFile string) bool {
	//fmt.Println("Downloader Start...")

	err := ipfs.Export2CarFileByIpfsUrl(downUrl, tarFile)
	if err != nil {
		logs.GetLogger().Error("Downloader Export2CarFileByIpfs :", err)
		return false
	}

	logs.GetLogger().Info("Downloader Result:", "FromUrl=", downUrl, " tarFile=", tarFile)

	return true
}

//url  "http://127.0.0.1:5001"
//fileHash  "QmTKZA24mw8gNTLnZeTZC6maSaHoDy8VB4eWBbtKF4gF1Z"
func CheckIpfsAlive(url string, fileHash string) bool {
	r := false

	cl := goipfs.NewShell(url)
	_, err := cl.List(fileHash)

	if err == nil {
		r = true
	}

	return r
}
