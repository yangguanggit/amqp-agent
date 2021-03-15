package config

import (
	"amqp-agent/common/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Project  string   `yaml:"project"`
	Database Database `yaml:"database"`
	Mq       Mq       `yaml:"mq"`
}

type Database struct {
	Dialect        string
	Host           string
	Port           int
	Database       string
	User           string
	Password       string
	Charset        string
	MaxIdleConnNum int
	MaxOpenConnNum int
}

type Mq struct {
	Host     string
	Port     int
	User     string
	Password string
	Min      int
	Max      int
}

var (
	application string
	runEnv      string
	AppConfig   *Config
	once        sync.Once
)

/**
 * 初始化配置
 */
func InitConfig(app, env, path string) {
	once.Do(func() {
		application = app
		runEnv = env
		if path == "" {
			pwd, _ := os.Getwd()
			path = filepath.Join(pwd, "config", env) + ".yaml"
		}
		AppConfig = loadConfig(path)
		//日志级别
		level := uint32(6)
		if env == gin.ReleaseMode {
			level = uint32(4)
		}
		logger.InitLogger(AppConfig.Project, application, level)
		gin.SetMode(env)
	})
}

/**
 * 获取应用名称
 */
func GetApplication() string {
	return application
}

/**
 * 获取运行时环境
 */
func GetEnv() string {
	return runEnv
}

/**
 * 加载配置
 */
func loadConfig(path string) *Config {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("fatal error: read config file: %s\n", err)
	}

	config := new(Config)
	if err := yaml.Unmarshal(content, config); err != nil {
		log.Fatalf("fatal error: load config file: %s\n", err)
	}

	return config
}
