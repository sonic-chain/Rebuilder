package mcs

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type BucketInfo struct {
	Name       string `json:"name"`
	FileNumber int    `json:"file_number"`
	Size       int    `json:"size"`
	MaxSize    int64  `json:"max_size"`
}

func GetBuckets(ctx context.Context) ([]BucketInfo, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	currentDir, _ := os.Getwd()
	cmd := exec.CommandContext(ctx, "python", path.Join(currentDir, "sdk.py"), "bucket_list")
	cmd.Stderr = cmd.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "execute cmd failed")
	}

	r := bufio.NewReader(stdout)
	line, err := readLine(r, "[{'")
	if err != nil {
		return nil, err
	}

	var infos []BucketInfo
	line = strings.ReplaceAll(line, "'", "\"")
	if line != "" {
		if err := json.Unmarshal([]byte(line), &infos); err != nil {
			return nil, errors.Wrap(err, "json Unmarshal failed")
		}
	}
	return infos, nil
}

func UploadFile(ctx context.Context, bucketName, objectName, localFileName string) ([]byte, error) {
	currentDir, _ := os.Getwd()
	uploadCmd := exec.CommandContext(ctx, "python", path.Join(currentDir, "sdk.py"), "upload_file", bucketName, objectName, localFileName)
	return uploadCmd.CombinedOutput()
}

func GetFile(ctx context.Context, bucketName, objectName string) (string, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	currentDir, _ := os.Getwd()
	cmd := exec.CommandContext(ctx, "python", path.Join(currentDir, "sdk.py"), "get_file", bucketName, objectName)
	cmd.Stderr = cmd.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err = cmd.Start(); err != nil {
		return "", errors.Wrap(err, "execute cmd failed")
	}

	r := bufio.NewReader(stdout)
	var ipfsUrl string

	line, err := readLine(r, "ipfs_url")
	if err != nil {
		return "", err
	}

	if line != "" {
		split := strings.Split(line, ": ")
		ipfsUrl = split[1][1 : len(split[1])-3]
	}
	return ipfsUrl, nil
}

func readLine(r *bufio.Reader, condition string) (string, error) {
	var data string
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", errors.Wrap(err, "read stdout failed")
			}
		}

		if strings.Contains(line, condition) {
			data = line
			break
		}
	}
	return data, nil
}
