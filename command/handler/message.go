package handler

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/logger"
	"amqp-agent/common/model"
	"amqp-agent/common/util"
	"amqp-agent/service"
	"amqp-agent/web/form"
	"context"
	"encoding/json"
)

type DelayMessage struct{}

//处理延时消息
func (a *DelayMessage) Handle(ctx context.Context) {
	delayList, err := new(model.DelayMessage).GetList()
	if err != nil {
		return
	}
	for _, delay := range delayList {
		if delay.DelayMessageStatus != constant.MessageStatusInit {
			continue
		}

		var err error
		var data *service.Result
		switch delay.DelayMessageType {
		case constant.MessageTypeQueue:
			params := new(form.QueueSendForm)
			if err := json.Unmarshal([]byte(delay.DelayMessageData), params); err != nil {
				logger.Error(ctx, "延时队列处理异常 - 参数解析错误", map[string]interface{}{
					"id":    delay.DelayMessageId,
					"error": err.Error(),
				}, constant.CateDelay)
				continue
			}
			messageList, _ := util.ParseMessage(params.Data)
			//队列消息
			data, err = new(service.Producer).EnQueue(ctx, params.QueueName, messageList)
		case constant.MessageTypeTopic:
			params := new(form.TopicPublishForm)
			if err := json.Unmarshal([]byte(delay.DelayMessageData), params); err != nil {
				logger.Error(ctx, "延时队列处理异常 - 参数解析错误", map[string]interface{}{
					"id":    delay.DelayMessageId,
					"error": err.Error(),
				}, constant.CateDelay)
				continue
			}
			messageList, _ := util.ParseMessage(params.Data)
			//主题消息
			data, err = new(service.Producer).PublishTopic(ctx, params.TopicName, params.RoutingKey, params.TagList, messageList)
		}

		//更新记录
		if data != nil {
			js, _ := json.Marshal(data)
			_ = delay.UpdateStatus(constant.MessageStatusSuccess, string(js))
		} else if err != nil {
			logger.Error(ctx, "延时队列处理异常 - 发送消息错误", map[string]interface{}{
				"id":    delay.DelayMessageId,
				"error": err.Error(),
			}, constant.CateDelay)
			_ = delay.UpdateStatus(constant.MessageStatusFail, err.Error())
		}
	}
}
