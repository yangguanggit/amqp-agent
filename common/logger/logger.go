package logger

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	project     string
	application string
	logger      *logrus.Logger
	loggerOnce  sync.Once
)

/**
 * 初始化日志组件
 */
func InitLogger(projectName string, applicationName string, level uint32) {
	loggerOnce.Do(func() {
		project, application = projectName, applicationName
		logger = logrus.New()
		logger.Level = logrus.Level(level)
		//控制台打印即可
		logger.SetOutput(os.Stdout)
		//json格式输出，PrettyPrint不能使用，否则k8s采集的有问题
		logger.SetFormatter(&logrus.JSONFormatter{})
	})
}

/**
 * trace
 */
func Trace(ctx context.Context, body ...interface{}) {
	entry, message := getEntry(ctx, constant.CateTrace, constant.CateTrace), formatMessage(body)
	entry.Trace(message)
}

/**
 * info
 */
func Info(ctx context.Context, overview string, body interface{}, category string) {
	entry, message := getEntry(ctx, overview, category), formatMessage(body)
	entry.Info(message)
}

/**
 * warning
 */
func Warning(ctx context.Context, overview string, body interface{}, category string) {
	entry, message := getEntry(ctx, overview, category), formatMessage(body)
	entry.Warning(message)
}

/**
 * error
 */
func Error(ctx context.Context, overview string, body interface{}, category string) {
	entry, message := getEntry(ctx, overview, category), formatMessage(body)
	entry.Error(message)
}

/**
 * 获取日志句柄
 */
func getEntry(ctx context.Context, overview string, category string) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"project":     project,
		"application": application,
		"traceId":     util.GetTraceId(ctx),
		"overview":    overview,
		"category":    category,
	})
}

/**
 * 格式化消息
 */
func formatMessage(body interface{}) string {
	data := ""
	switch value := body.(type) {
	case string:
		data = value
	case error:
		data = value.Error()
	default:
		if logData, err := json.MarshalIndent(body, "", "    "); err == nil {
			data = string(logData)
		} else {
			data = fmt.Sprint(body)
		}
	}
	if len(data) > 8192 {
		return data[0:8192]
	}
	return data
}
