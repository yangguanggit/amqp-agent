package mq

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/logger"
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidConfig = errors.New("MQ: Invalid pool config")
	ErrPoolClosed    = errors.New("MQ: Pool closed")
)

type factory func() (*MQ, error)

type Pool interface {
	Get() (*MQ, error) // 获取资源
	Put(*MQ) error     // 释放资源
	Close(*MQ) error   // 关闭资源
	ClosePool() error  // 关闭池
}

type GenericPool struct {
	sync.Mutex
	pool        chan *MQ
	minOpen     int  // 池中最少资源数
	maxOpen     int  // 池中最大资源数
	numOpen     int  // 当前池中资源数
	closed      bool // 池是否已关闭
	maxLifetime time.Duration
	factory     factory // 创建连接的方法
}

func NewGenericPool(minOpen, maxOpen int, factory factory) (*GenericPool, error) {
	if minOpen <= 0 || maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}
	p := &GenericPool{
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		maxLifetime: 0,
		factory:     factory,
		pool:        make(chan *MQ, maxOpen),
	}

	for i := 0; i < minOpen; i++ {
		resource, err := factory()
		if err != nil {
			logger.Error(nil, "Init generic pool error", err, constant.CateAmqp)
			continue
		}
		p.numOpen++
		p.pool <- resource
	}
	return p, nil
}

func (p *GenericPool) Get() (*MQ, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}
	resource, err := p.getOrCreate()
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (p *GenericPool) getOrCreate() (*MQ, error) {
	select {
	case resource := <-p.pool:
		return resource, nil
	default:
	}
	p.Lock()
	if p.numOpen >= p.maxOpen {
		p.Unlock()
		resource := <-p.pool
		return resource, nil
	}
	p.numOpen++
	p.Unlock()

	// 新建连接
	resource, err := p.factory()
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// 释放单个资源到连接池
func (p *GenericPool) Put(resource *MQ) error {
	if p.closed {
		return ErrPoolClosed
	}
	p.pool <- resource
	return nil
}

// 关闭单个资源
func (p *GenericPool) Close(resource *MQ) error {
	p.Lock()
	p.numOpen--
	p.Unlock()
	resource.Close()
	return nil
}

// 关闭连接池，释放所有资源
func (p *GenericPool) ClosePool() error {
	if p.closed {
		return nil
	}
	p.Lock()
	defer p.Unlock()
	close(p.pool)
	for resource := range p.pool {
		resource.Close()
		p.numOpen--
	}
	p.closed = true
	return nil
}
