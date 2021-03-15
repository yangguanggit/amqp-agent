# AMQP-AGENT

### 说明
##### 提供`amqp`协议消息中间件的`http`代理服务

### 准备
##### 1.部署`rabbitmq`，后台创建`direct`类型`exchange` - `exchange.queue`
##### 2.部署`mysql`，初始化数据库`common/model/db.sql`
##### 3.修改`config`目录下不同环境的配置文件

### 启动
```sh
cd amqp-agent
go mod download
# web服务
go run web/main.go
# 脚本服务
go run command/main.go
```

### 部署
```sh
# 编译
bash build.sh release web main
# 启动web应用
bash start.sh release main
# 平滑重启web应用
bash reload.sh release main
```

### 接口列表
| path | method | 说明 |
| ---- | ---- | ---- |
| /amqp/queue/pull | GET | 拉取消息 |
| /amqp/queue/ack | POST | 确认消息（暂未实现拉取、确认分离） |
| /amqp/queue/send | POST | 发送队列消息 |
| /amqp/topic/publish | POST | 发送主题消息 |
