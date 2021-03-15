package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type Consumer struct {
	name     string             // Consumer的名字
	state    uint8              // Consumer状态
	prefetch int                // Qos prefetch
	callback chan amqp.Delivery // 上层用于接收消费出来的消息的管道
	ch       *amqp.Channel      // MQ的会话channel
}

func NewConsumer(name string) *Consumer {
	return &Consumer{
		name: name,
	}
}

// 设置channel每次消费消息数量, prefetch默认为0取值范围[0,∞)
func (c *Consumer) SetQos(prefetch int) *Consumer {
	c.prefetch = prefetch
	return c
}

func (c *Consumer) SetCallback(callback chan amqp.Delivery) *Consumer {
	c.callback = callback
	return c
}

func (c *Consumer) Open() error {
	// 状态检测
	if c.state == StateOpened {
		return fmt.Errorf("MQ: Consumer(%s) had been opened", c.name)
	}
	// 参数校验
	if c.name == "" {
		return fmt.Errorf("MQ: Consumer(%s) empty name", c.name)
	}

	// 初始化channel
	ch, err := GetMQConnect()
	if err != nil {
		return err
	}
	c.ch = ch
	// 检查队列
	queue := DefaultQueue(c.name)
	if _, err := c.ch.QueueDeclarePassive(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, queue.NoWait, queue.Args); err != nil {
		return fmt.Errorf("MQ: Declare queue(%s) failed, %v", queue.Name, err)
	}
	// 设置属性
	if err := c.ch.Qos(c.prefetch, 0, false); err != nil {
		return fmt.Errorf("MQ: Consumer(%s) open error, %v", c.name, err)
	}
	c.state = StateOpened

	return nil
}

func (c *Consumer) Close() {
	if c.ch != nil {
		_ = c.ch.Close()
	}
	c.state = StateClosed
}

func (c *Consumer) Consume(opt *ConsumeOption) error {
	if c.state != StateOpened {
		return fmt.Errorf("MQ: Consumer(%s) unopened, now state is %d", c.name, c.state)
	}
	if opt.MessageCount <= 0 || opt.BlockSecond <= 0 {
		return fmt.Errorf("MQ: Consumer(%s) option is invalid", c.name)
	}
	if c.callback == nil {
		return fmt.Errorf("MQ: Consumer(%s) callback is nil", c.name)
	}

	delivery, err := c.ch.Consume(c.name, "", opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args)
	if err != nil {
		return fmt.Errorf("MQ: Consumer(%s) consume error, %v", c.name, err)
	}
	go c.deliver(opt, delivery)

	return nil
}

func (c *Consumer) deliver(opt *ConsumeOption, delivery <-chan amqp.Delivery) {
	ch := time.After(time.Duration(opt.BlockSecond) * time.Second)

	for i := 1; i <= opt.MessageCount; i++ {
		select {
		case d := <-delivery:
			c.callback <- d
			if i >= opt.MessageCount {
				//消息数量达到关闭通道
				close(c.callback)
				return
			}
		case <-ch:
			//超时关闭通道
			close(c.callback)
			return
		}
	}
}
