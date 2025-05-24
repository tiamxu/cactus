package types

import (
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/pkg/e"
)

type BasePage struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

// DataListResp 带有总数的Data结构

type DataListResp struct {
	PageData interface{} `json:"PageData"`
	Total    int64       `json:"total"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Meta    interface{} `json:"meta,omitempty"`
	// OriginUrl string      `json:"originUrl"`
}

func RespSuccess(c *gin.Context, data interface{}, code ...int) *Response {
	status := e.SUCCESS
	if code != nil {
		status = code[0]
	}

	if data == nil {
		data = "操作成功"
	}

	r := &Response{
		Code:    status,
		Data:    data,
		Message: e.GetMsg(status),
		// OriginUrl: c.Request.URL.Path,
	}
	return r
}

func RespError(c *gin.Context, err error, data string, code ...int) *Response {
	status := e.ERROR
	if code != nil {
		status = code[0]
	}
	r := &Response{
		Code:    status,
		Data:    data,
		Message: e.GetMsg(status),
		Error:   err.Error(),
	}
	return r
}
