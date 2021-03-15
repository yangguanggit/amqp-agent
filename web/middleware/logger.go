package middleware

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/logger"
	"amqp-agent/common/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 日志记录到控制台
func LoggerToStdout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 请求方式
		reqMethod := c.Request.Method
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()

		// 请求返回值
		var request, response []byte
		if requestData, ok := c.Get(gin.BodyBytesKey); ok && requestData != nil {
			if requestBytes, ok := requestData.([]byte); ok {
				request = requestBytes
			}
		}
		if responseData, ok := c.Get(constant.SysResponseDataKey); ok && responseData != nil {
			response, _ = json.Marshal(responseData)
		}

		requestURI := c.Request.RequestURI
		//如果uri中有参数，则截取掉
		if c.FullPath() != "" {
			requestURI = strings.Split(c.FullPath(), "/:")[0]
		}

		data := map[string]interface{}{
			"method":   reqMethod,
			"uri":      requestURI,
			"ua":       c.Request.UserAgent(),
			"clientIp": clientIP,
			"status":   statusCode,
			"spend":    float64(endTime.Sub(startTime)) / 1e6,
			"request":  string(request),
			"response": string(response),
		}

		// 日志记录
		ctx, _ := util.ContextWithSpan(c)
		logger.Info(ctx, "GIN-INFO", data, constant.CateTrace)
	}
}
