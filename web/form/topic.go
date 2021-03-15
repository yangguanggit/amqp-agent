package form

type TopicPublishForm struct {
	Source     string      `json:"source"`
	TopicName  string      `json:"topicName" binding:"required"`
	Data       interface{} `json:"data" binding:"required"`
	TagList    []string    `json:"tagList"`
	RoutingKey string      `json:"routingKey"`
	DelayTime  int         `json:"delayTime"`
}
