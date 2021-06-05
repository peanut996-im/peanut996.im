// Package api
// @Title  api_response.go
// @Description  record some defined response.
// @Author  peanut996
// @Update  peanut996  2021/5/22 0:22
package api

import (
	"fmt"
)

type BaseRepsonse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	SignInvaildResp = &BaseRepsonse{
		Code:    ErrorSignInvalid,
		Message:  ErrorCodeToString(ErrorSignInvalid),
		Data:    nil,
	}

	ResourceExistsResp = &BaseRepsonse{
		Code:    ErrorHttpResourceExists,
		Message: ErrorCodeToString(ErrorHttpResourceExists),
		Data:    nil,
	}

	ResourceNotFoundResp = &BaseRepsonse{
		Code:    ErrorHttpResourceNotFound,
		Message: ErrorCodeToString(ErrorHttpResourceNotFound),
		Data:    nil,
	}

	AuthFaildResp = &BaseRepsonse{
		Code:    ErrorAuthFailed,
		Message: ErrorCodeToString(ErrorAuthFailed),
		Data:    nil,
	}

	TokenInvaildResp = &BaseRepsonse{
		Code:    ErrorTokenInvalid,
		Message: ErrorCodeToString(ErrorTokenInvalid),
		Data:    nil,
	}
)

func NewBaseResponse(code int, data interface{}) *BaseRepsonse {
	return &BaseRepsonse{
		Code:    code,
		Data:    data,
		Message: ErrorCodeToString(code),
	}
}

func NewHttpInnerErrorResponse(err error) *BaseRepsonse {
	return &BaseRepsonse{
		Code:    ErrorHttpInnerError,
		Message: fmt.Sprintf(ErrorCodeToFormat(ErrorHttpInnerError), err),
		Data:    nil,
	}
}

func NewSuccessResponse(data interface{}) *BaseRepsonse {
	return &BaseRepsonse{
		Code:    ErrorCodeOK,
		Data:    data,
		Message: ErrorCodeToFormat(ErrorCodeOK),
	}
}

func NewResourceExistsResponse(err error) *BaseRepsonse {
	return &BaseRepsonse{
		Code:    ErrorHttpResourceExists,
		Data:    nil,
		Message: fmt.Sprintf(ErrorCodeToFormat(ErrorHttpResourceExists), err),
	}
}
