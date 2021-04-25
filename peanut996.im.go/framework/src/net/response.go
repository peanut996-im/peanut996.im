package net

const (
	ERROR_CODE_OK = iota
	ERROR_TOKEN_INVALID
	ERROR_HTTP_INNER_ERROR
	ERROR_HTTP_PARAM_INVALID
	ERROR_REDIS_ERROR
	ERROR_MONGO_ERROR
)

var (
	errorCodeInfo = map[int]string{
		ERROR_CODE_OK:            "success",
		ERROR_TOKEN_INVALID:      "token invalid",
		ERROR_HTTP_INNER_ERROR:   "http inner error",
		ERROR_HTTP_PARAM_INVALID: "http param invalid",
		ERROR_REDIS_ERROR:        "redis error",
		ERROR_MONGO_ERROR:        "mongo error",
	}
)

type BaseRepsonse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ErrorCodeToString(code int) string {
	return errorCodeInfo[code]
}
