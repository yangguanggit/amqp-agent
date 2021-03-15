package service

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/logger"
	"amqp-agent/common/mq"
	"context"
	"github.com/streadway/amqp"
	"time"
)

type Consumer struct{}

type Message struct {
	MessageId   string `json:"messageId"`
	MessageTag  uint64 `json:"messageTag"`
	MessageBody string `json:"messageBody"`
	EnqueueTime int64  `json:"enqueueTime"`
	DequeueTime int64  `json:"dequeueTime"`
}

//消息出队列
func (c *Consumer) DeQueue(ctx context.Context, queueName string, messageCount, blockSecond int) ([]*Message, error) {
	consumer := mq.NewConsumer(queueName)
	deliverChan := make(chan amqp.Delivery, messageCount)
	defer func() {
		consumer.Close()
		//由于内部达到消息数量或空消息时间限制会关闭通道，关闭前检查一下状态
		if _, ok := <-deliverChan; ok {
			close(deliverChan)
		}
	}()
	if err := consumer.SetQos(messageCount).SetCallback(deliverChan).Open(); err != nil {
		close(deliverChan)
		return nil, err
	}

	// 开始循环消费
	opt := mq.DefaultConsumeOption(messageCount, blockSecond)
	if err := consumer.Consume(opt); err != nil {
		close(deliverChan)
		return nil, err
	}

	messageList := make([]*Message, 0)
	for d := range deliverChan {
		message := &Message{
			MessageId:   d.MessageId,
			MessageTag:  d.DeliveryTag,
			MessageBody: string(d.Body),
			EnqueueTime: d.Timestamp.Unix(),
			DequeueTime: time.Now().Unix(),
		}
		messageList = append(messageList, message)
		if err := d.Ack(false); err != nil {
			logger.Error(ctx, "消息确认异常", err.Error(), constant.CateConsumer)
		}
	}

	return messageList, nil
}

//确认消息
func (c *Consumer) Ack(ctx context.Context, queueName string, messageIdList []string) error {
	//todo pull、ack分离
	return nil
}
