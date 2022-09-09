package ipfs

import (
	"fmt"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func IpfsUploadFileByWebApi(apiUrl, filefullpath string) (*string, error) {
	response, err := web.HttpUploadFileByStream(apiUrl, filefullpath)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	fileHash := utils.GetFieldStrFromJson(response, "Hash")

	if fileHash == constants.EMPTY_STRING {
		err := fmt.Errorf("cannot get file hash from response:%s", response)
		//logs.GetLogger().Error(err)
		return nil, err
	}

	return &fileHash, nil
}

func IpfsUploadFileByStream(apiUrl, filefullpath string) (string, error) {
	response, err := web.HttpUploadFileByStream(apiUrl, filefullpath)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	fileHash := utils.GetFieldStrFromJson(response, "Hash")

	if fileHash == constants.EMPTY_STRING {
		err := fmt.Errorf("cannot get file hash from response:%s", response)
		//logs.GetLogger().Error(err)
		return "", err
	}

	return fileHash, nil
}

func Export2CarFile(apiUrl, fileHash string, carFileFullPath string) error {
	apiUrlFull := utils.UrlJoin(apiUrl, "api/v0/dag/export")
	apiUrlFull = apiUrlFull + "?arg=" + fileHash + "&progress=false"
	fmt.Println("apiUrlFull :", apiUrlFull)

	carFileContent, err := web.HttpPostNoToken(apiUrlFull, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	fmt.Println("Export2CarFile :", string(carFileContent))

	bytesWritten, err := utils.CreateFileWithByteContents(carFileFullPath, carFileContent)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	logs.GetLogger().Info(bytesWritten, " bytes have been written to:", carFileFullPath)
	return nil
}

func Export2CarFileByIpfs(apiUrl, fileHash string, carFileFullPath string) error {
	apiUrlFull := utils.UrlJoin(apiUrl, "ipfs/"+fileHash)

	carFileContent, err := web.HttpGetNoToken(apiUrlFull, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	bytesWritten, err := utils.CreateFileWithByteContents(carFileFullPath, carFileContent)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	logs.GetLogger().Info(bytesWritten, " bytes have been written to:", carFileFullPath)
	return nil
}

func Export2CarFileByIpfsUrl(ipfsFullUrl, carFileFullPath string) error {

	carFileContent, err := web.HttpGetNoToken(ipfsFullUrl, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	_, err = utils.CreateFileWithByteContents(carFileFullPath, carFileContent)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	//logs.GetLogger().Info(bytesWritten, " bytes have been written to:", carFileFullPath)
	return nil
}
