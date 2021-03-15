package main

import (
	"amqp-agent/command/cmd"
	"amqp-agent/common/constant"
	"amqp-agent/common/logger"
	"amqp-agent/common/model"
	"amqp-agent/common/mq"
	"amqp-agent/common/util"
	"amqp-agent/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
)

func init() {
	//初始化配置文件
	env := util.GetArg(constant.ArgEnv, gin.DebugMode)
	path := util.GetArg(constant.ArgConfig, "")
	config.InitConfig(constant.SysCommand, env, path)
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			//打印调用栈信息
			s := make([]byte, 2048)
			n := runtime.Stack(s, false)
			logger.Error(nil, "command exception", fmt.Sprintf("%s, %s", p, s[:n]), constant.CateTrace)
		}
	}()

	//注册数据库连接
	model.ConnectDB(config.AppConfig.Database)
	defer model.CloseDB()

	//连接mq
	if err := mq.InitMqPool(config.AppConfig.Mq); err != nil {
		log.Fatalf("mq connect error: %s\n", err)
	}
	defer mq.CloseMqPool()

	cmd.Execute()
}
