package util

import (
	"amqp-agent/common/constant"
	"context"
	"github.com/gin-gonic/gin"
)

// ContextWithSpan 返回context
func ContextWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(constant.SysContextTracerKey)
	if !exist {
		ctx, ok = context.TODO(), false
		return
	}

	ctx, ok = v.(context.Context)
	return
}

func ContextWithCancel(c *gin.Context) (ctx context.Context, cancel context.CancelFunc, ok bool) {
	ctx, ok = ContextWithSpan(c)
	ctx, cancel = context.WithCancel(ctx)
	return
}

/**
 * 获取链路tranceId
 */
func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	//todo tranceId
	return ""
}
