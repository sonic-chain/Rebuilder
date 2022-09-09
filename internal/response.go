package internal

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_FILE_LIST_FAIL = 6001
	ERROR_CHANGETO_JSON  = 6002
	ERROR_MINERID_FAIL   = 6003
	ERROR_RETRIEVE_FAIL  = 6004
	ERROR_UPLOAD_FAIL    = 6005
	ERROR_SUMMARY_FAIL   = 6006
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})
	return
}

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Info(err.Key, err.Message)
	}

	return
}
