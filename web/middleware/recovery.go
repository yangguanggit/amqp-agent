package middleware

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/logger"
	"amqp-agent/common/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				//打印调用栈信息
				s := make([]byte, 2048)
				n := runtime.Stack(s, false)

				ctx, _ := util.ContextWithSpan(c)
				logger.Error(ctx, "GIN-ERROR", fmt.Sprintf("%s, %s", err, s[:n]), constant.CateTrace)

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		c.Next()
	}
}
