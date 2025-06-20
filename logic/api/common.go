package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/pkg/e"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
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

func RespError(c *gin.Context, err error, msg string, code ...int) *Response {
	status := e.ERROR
	if code != nil {
		status = code[0]
	}
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	r := &Response{
		Code:    status,
		Data:    "",
		Message: msg,
		Error:   errorMsg,
	}
	return r
}

/*
// 成功响应
{
    "code": 200,
    "message": "操作成功",
    "data": { ... }  // 实际业务数据
}

// 错误响应
{
    "code": 400,
    "message": "验证码错误",       // 用户友好提示
    "error": "captcha_mismatch", // 错误代码（可选，用于前端逻辑）
    "error_details": "got '1234', expected '5678'" // 可选，开发阶段可见
}
*/
