package util

import (
	"amqp-agent/common/constant"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * 成功
 */
func Success(ctx *gin.Context, data interface{}) {
	d := gin.H{
		"code":      constant.ErrorOk,
		"message":   constant.ErrorMessage(constant.ErrorOk),
		"data":      data,
		"requestId": "",
	}
	ctx.Set(constant.SysResponseDataKey, d)
	ctx.JSON(http.StatusOK, d)
}

/**
 * 失败
 */
func Fail(ctx *gin.Context, errs ...interface{}) {
	code := constant.ErrorSystemError
	message := constant.ErrorMessage(code)
	for _, e := range errs {
		switch value := e.(type) {
		case int:
			code = value
			message = constant.ErrorMessage(code)
		case string:
			message = value
		case error:
			message = value.Error()
		}
	}
	data := gin.H{
		"code":      code,
		"message":   message,
		"data":      gin.H{},
		"requestId": "",
	}
	ctx.Set(constant.SysResponseDataKey, data)
	ctx.JSON(http.StatusOK, data)
}

/**
 * 错误
 * 必须:错误编号
 * 可选:系统错误，自定义业务错误信息
 */
func Error(code int, errs ...interface{}) error {
	err := errors.New(constant.ErrorMessage(code))
	for _, e := range errs {
		switch value := e.(type) {
		case string:
			err = errors.New(value)
		case error:
			err = value
		}
	}

	return err
}
