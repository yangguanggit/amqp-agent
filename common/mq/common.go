package mq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

const (
	StateClosed    = uint8(0)
	StateOpened    = uint8(1)
	StateReopening = uint8(2)
)

// ExchangeBind exchange => routeKey => queue
type ExchangeBind struct {
	Exchange *Exchange
	Binding  *Binding
}

// Exchange配置
type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table // default is nil
}

// Biding routeKey => queue
type Binding struct {
	RouteKey string
	Queue    *Queue
	NoWait   bool       // default is false
	Args     amqp.Table // default is nil
}

// Queue配置
type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// 消费者消费选项
type ConsumeOption struct {
	AutoAck      bool
	Exclusive    bool
	NoLocal      bool
	NoWait       bool
	Args         amqp.Table
	MessageCount int //每次拉取消息数量
	BlockSecond  int //空消息时等待秒数
}

func DefaultExchange(name string, kind string) *Exchange {
	return &Exchange{
		Name:       name,
		Kind:       kind,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

func DefaultQueue(name string) *Queue {
	return &Queue{
		Name:       name,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

func NewPublishMessage(body []byte, headers amqp.Table) amqp.Publishing {
	return amqp.Publishing{
		UserId:       "", //todo 用户标记
		AppId:        "", //todo 应用标记
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Priority:     uint8(5),
		Timestamp:    time.Now(),
		MessageId:    strconv.FormatInt(time.Now().UnixNano(), 10),
		Headers:      headers,
		Body:         body,
	}
}

func DefaultConsumeOption(messageCount, blockSecond int) *ConsumeOption {
	return &ConsumeOption{
		NoWait:       true,
		MessageCount: messageCount,
		BlockSecond:  blockSecond,
	}
}

func applyExchangeBind(ch *amqp.Channel, exchangeBind *ExchangeBind) error {
	if ch == nil {
		return errors.New("MQ: Nil channel")
	}
	if exchangeBind == nil {
		return errors.New("MQ: Nil exchangeBind")
	}
	if exchangeBind.Exchange == nil {
		return errors.New("MQ: Nil exchange")
	}
	// todo 主题exchange和queue绑定关系后台配置，队列exchange和queue绑定关系才需要声明
	/*if exchangeBind.Binding == nil {
		return fmt.Errorf("MQ: Nil binding")
	}*/

	// declare exchange
	exchange := exchangeBind.Exchange
	// todo exchange和queue信息后台手动添加配置，没有配置直接返回错误信息
	if err := ch.ExchangeDeclarePassive(exchange.Name, exchange.Kind, exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, exchange.Args); err != nil {
		return fmt.Errorf("MQ: Declare exchange(%s) failed, %v", exchange.Name, err)
	}
	// declare queue and bind queue
	if binding := exchangeBind.Binding; binding != nil {
		if queue := binding.Queue; queue != nil {
			// todo exchange和queue信息后台手动添加配置，没有配置直接返回错误信息
			if _, err := ch.QueueDeclarePassive(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, queue.NoWait, queue.Args); err != nil {
				return fmt.Errorf("MQ: Declare queue(%s) failed, %v", queue.Name, err)
			}
			if err := ch.QueueBind(queue.Name, binding.RouteKey, exchange.Name, binding.NoWait, binding.Args); err != nil {
				return fmt.Errorf("MQ: Bind exchange(%s) <-> queue(%s) failed, %v", exchange.Name, queue.Name, err)
			}
		}
	}
	return nil
}
