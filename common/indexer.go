package common

import (
	"crypto/tls"
	"github.com/Fogmeta/filecoin-ipfs-data-rebuilder/model"
	"io"
	"strings"
	"sync"

	"net/http"
	"time"
)

const (
	GET_PEER_URL = "/cid/"
)

type indexer struct {
	urls   []string
	client *http.Client
}

func NewIndexerClient() *indexer {
	if len(model.IndexerSetting.Urls) == 0 || model.IndexerSetting.Urls[0] == "" {
		log.Fatalf("read indexer config urls failed, please check indexer.urls.")
	}
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	return &indexer{
		urls:   model.IndexerSetting.Urls,
		client: &client,
	}
}

func (indexer *indexer) SendHttpGet(url, cid string) (result [][]byte) {
	result = make([][]byte, 0)
	var wg sync.WaitGroup
	wg.Add(len(indexer.urls))
	for i := 0; i < len(indexer.urls); i++ {
		num := i
		go func() {
			defer wg.Done()
			reqUrl := indexer.urls[num] + url + cid
			log.Infof("send indexer node url: %s", reqUrl)
			request, err := http.NewRequest("GET", reqUrl, nil)
			if err != nil {
				log.Errorf("create get cid request failed,error:%+v", err)
				return
			}
			resp, err := indexer.client.Do(request)
			if err != nil {
				log.Errorf("send cid request failed,error:%+v", err)
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("read cid request failed,error:%+v", err)
				return
			}
			if !strings.EqualFold(string(data), "no results for query") {
				result = append(result, data)
			}
		}()
	}
	wg.Wait()
	return result
}
