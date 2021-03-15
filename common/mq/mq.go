package mq

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/logger"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type MQ struct {
	sync.Mutex                  // 保护内部数据并发读写
	addr       string           // MQ连接的地址
	state      uint8            // MQ状态
	conn       *amqp.Connection // MQ TCP连接
	closeChan  chan *amqp.Error // MQ监听连接错误
	stopChan   chan struct{}    // 监听用户主动关闭
}

func NewMQ(addr string) *MQ {
	return &MQ{
		addr:  addr,
		state: StateClosed,
	}
}

func (m *MQ) Open() (mq *MQ, err error) {
	// 进行Open期间不允许做任何跟连接有关的事情
	m.Lock()
	defer m.Unlock()

	if m.state == StateOpened {
		return m, errors.New("MQ: Had been opened")
	}

	if m.conn, err = m.dial(); err != nil {
		return m, fmt.Errorf("MQ: Dial err: %v", err)
	}

	m.state = StateOpened
	m.stopChan = make(chan struct{})
	m.closeChan = make(chan *amqp.Error, 1)
	m.conn.NotifyClose(m.closeChan)

	go m.keepalive()

	return m, nil
}

func (m *MQ) Close() {
	m.Lock()
	// close mq connection
	select {
	case <-m.stopChan:
		// had been closed
	default:
		close(m.stopChan)
	}
	m.Unlock()

	// wait done
	for m.State() != StateClosed {
		time.Sleep(time.Second)
	}
}

func (m *MQ) State() uint8 {
	m.Lock()
	defer m.Unlock()
	return m.state
}

func (m MQ) dial() (*amqp.Connection, error) {
	return amqp.Dial(m.addr)
}

func (m *MQ) Channel() (*amqp.Channel, error) {
	return m.conn.Channel()
}

func (m *MQ) keepalive() {
	select {
	case <-m.stopChan:
		// 正常关闭
		logger.Info(nil, "MQ: Close normally", nil, constant.CateAmqp)
		m.Lock()
		_ = m.conn.Close()
		m.state = StateClosed
		m.Unlock()

	case err := <-m.closeChan:
		logger.Error(nil, "MQ: Disconnected with MQ", err, constant.CateAmqp)
		// tcp连接中断, 重新连接
		m.Lock()
		m.state = StateReopening
		m.Unlock()

		maxRetry := 10000
		for i := 1; i <= maxRetry; i++ {
			time.Sleep(time.Second)
			if _, err := m.Open(); err != nil {
				logger.Error(nil, fmt.Sprintf("MQ: Connection recover failed for %d times", i), err, constant.CateAmqp)
				continue
			}
			logger.Info(nil, fmt.Sprintf("MQ: Connection recover OK for %d times", i), nil, constant.CateAmqp)
			return
		}
		logger.Error(nil, fmt.Sprintf("MQ: Connection recover failed over maxRetry(%d), so exit", maxRetry), nil, constant.CateAmqp)
	}
}
