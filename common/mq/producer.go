package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type Producer struct {
	name         string                 // Producer的名字
	state        uint8                  // Producer状态
	exchangeBind *ExchangeBind          // exchange与queue绑定关系
	confirm      bool                   // 生产者confirm开关
	confirmChan  chan amqp.Confirmation // 监听publish confirm
	ch           *amqp.Channel          // MQ的会话channel
}

func NewProducer(name string) (*Producer, error) {
	return &Producer{
		name: name,
	}, nil
}

// Confirm 是否开启生产者confirm功能, 默认为false, 该选项在Open()前设置.
// 说明: 目前仅实现串行化的confirm, 每次的等待confirm额外需要约50ms,建议上层并发调用Publish
func (p *Producer) SetConfirm(confirm bool) *Producer {
	p.confirm = confirm
	return p
}

func (p *Producer) SetExchangeBind(exchangeBind *ExchangeBind) *Producer {
	p.exchangeBind = exchangeBind
	return p
}

func (p *Producer) Open() error {
	// 状态检测
	if p.state == StateOpened {
		return fmt.Errorf("MQ: Producer(%s) had been opened", p.name)
	}
	// 参数校验
	if p.exchangeBind == nil {
		return fmt.Errorf("MQ: Producer(%s) nil exchangeBind, setExchangeBind before open", p.name)
	}

	// 初始化channel
	ch, err := GetMQConnect()
	if err != nil {
		return err
	}
	p.ch = ch
	// 检查主题、队列
	if err := applyExchangeBind(p.ch, p.exchangeBind); err != nil {
		return err
	}

	// 初始化发送Confirm
	if p.confirm {
		p.confirmChan = make(chan amqp.Confirmation, 1) // channel关闭时自动关闭
		_ = p.ch.Confirm(false)
		p.ch.NotifyPublish(p.confirmChan)
	}

	p.state = StateOpened

	return nil
}

func (p *Producer) Close() {
	if p.ch != nil {
		_ = p.ch.Close()
	}
	p.state = StateClosed
}

// 在同步Publish Confirm模式下, 每次Publish将额外有约50ms的等待时间. 如果采用这种模式, 建议上层并发publish
func (p *Producer) Publish(routeKey string, message amqp.Publishing) error {
	if p.state != StateOpened {
		return fmt.Errorf("MQ: Producer(%s) unopened, now state is %d", p.name, p.state)
	}

	// 非confirm模式
	if !p.confirm {
		return p.ch.Publish(p.exchangeBind.Exchange.Name, routeKey, false, false, message)
	}

	// confirm模式
	if err := p.ch.Publish(p.exchangeBind.Exchange.Name, routeKey, false, false, message); err != nil {
		return fmt.Errorf("MQ: Producer(%s) publish failed, %v", p.name, err)
	}
	select {
	case confirm, ok := <-p.confirmChan:
		if !ok || !confirm.Ack {
			return fmt.Errorf("MQ: Producer(%s) publish confirm failed, confirm:%v, ok:%t", p.name, confirm, ok)
		}
	case <-time.After(time.Second):
		return fmt.Errorf("MQ: Producer(%s) publish confirm timeout", p.name)
	}

	return nil
}
