package form

type QueueSendForm struct {
	Source    string      `json:"source"`
	QueueName string      `json:"queueName" binding:"required"`
	Data      interface{} `json:"data" binding:"required"`
	DelayTime int         `json:"delayTime"`
}

type QueueAckForm struct {
	QueueName string   `json:"queueName" binding:"required"`
	Data      []string `json:"data" binding:"required"`
}
