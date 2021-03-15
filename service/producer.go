package service

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/mq"
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

type Producer struct{}

type Result struct {
	Success []string `json:"success"`
	Fail    []string `json:"fail"`
}

//消息入队列
func (p *Producer) EnQueue(ctx context.Context, queueName string, messageList []string) (*Result, error) {
	exchangeBind := &mq.ExchangeBind{
		Exchange: mq.DefaultExchange(constant.ExchangeQueue, amqp.ExchangeDirect),
		Binding: &mq.Binding{
			RouteKey: queueName,
			Queue:    mq.DefaultQueue(queueName),
		},
	}

	producer, err := mq.NewProducer(queueName)
	if err != nil {
		return nil, err
	}
	defer producer.Close()
	if err := producer.SetExchangeBind(exchangeBind).SetConfirm(true).Open(); err != nil {
		return nil, err
	}
	result := new(Result)
	for _, d := range messageList {
		message := mq.NewPublishMessage([]byte(d), nil)
		if err := producer.Publish(queueName, message); err != nil {
			result.Fail = append(result.Fail, fmt.Sprintf("message:%v, error:%s", message, err.Error()))
			continue
		}
		result.Success = append(result.Success, message.MessageId)
	}

	return result, nil
}

//发送主题消息
func (p *Producer) PublishTopic(ctx context.Context, topicName, routingKey string, tagList, messageList []string) (*Result, error) {
	//路由模式 - Routing key: routingKey
	exchangeBind := &mq.ExchangeBind{
		Exchange: mq.DefaultExchange(topicName, amqp.ExchangeTopic),
	}
	//标签模式设置属性 - Arguments: tag=tag && x-match=any
	headers := amqp.Table{}
	if len(tagList) > 0 {
		for _, tag := range tagList {
			headers[tag] = tag
		}
		exchangeBind = &mq.ExchangeBind{
			Exchange: mq.DefaultExchange(topicName, amqp.ExchangeHeaders),
		}
	}

	producer, err := mq.NewProducer(topicName)
	if err != nil {
		return nil, err
	}
	defer producer.Close()
	if err := producer.SetExchangeBind(exchangeBind).SetConfirm(true).Open(); err != nil {
		return nil, err
	}
	result := new(Result)
	for _, d := range messageList {
		message := mq.NewPublishMessage([]byte(d), headers)
		if err := producer.Publish(routingKey, message); err != nil {
			result.Fail = append(result.Fail, fmt.Sprintf("message:%v, error:%s", message, err.Error()))
			continue
		}
		result.Success = append(result.Success, message.MessageId)
	}

	return result, nil
}
