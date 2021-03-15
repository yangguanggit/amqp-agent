package main

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/model"
	"amqp-agent/common/mq"
	"amqp-agent/common/util"
	"amqp-agent/config"
	"amqp-agent/web/middleware"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func init() {
	//初始化配置文件
	env := util.GetArg(constant.ArgEnv, gin.DebugMode)
	path := util.GetArg(constant.ArgConfig, "")
	config.InitConfig(constant.SysWeb, env, path)
}

func main() {
	//注册数据库连接
	model.ConnectDB(config.AppConfig.Database)
	defer model.CloseDB()

	//连接mq
	if err := mq.InitMqPool(config.AppConfig.Mq); err != nil {
		log.Fatalf("mq connect error: %s\n", err)
	}
	defer mq.CloseMqPool()

	engine := gin.New()
	engine.Use(
		middleware.GinRecovery(),
		middleware.LoggerToStdout(),
	)
	middleware.RegisterRoute(engine)
	//启动
	gracehttp.SetLogger(log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile))
	if err := gracehttp.Serve(&http.Server{Addr: ":80", Handler: engine}); err != nil {
		log.Fatalf("run service error: %s\n", err)
	}
}
