package internal

var MsgFlags = map[int]string{
	SUCCESS:              "ok",
	ERROR:                "fail",
	INVALID_PARAMS:       "Invalid request parameter error",
	ERROR_FILE_LIST_FAIL: "get file list error",
	ERROR_CHANGETO_JSON:  "convert json error",
	ERROR_MINERID_FAIL:   "get minerId error",
	ERROR_RETRIEVE_FAIL:  "Can not find any backup from Filecoin network",
	ERROR_UPLOAD_FAIL:    "upload file to IPFS error",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
