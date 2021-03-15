package mq

import (
	"amqp-agent/config"
	"fmt"
	"github.com/streadway/amqp"
)

var pool *GenericPool

//初始化连接池
func InitMqPool(mqConfig config.Mq) error {
	var err error
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d", mqConfig.User, mqConfig.Password, mqConfig.Host, mqConfig.Port)
	pool, err = NewGenericPool(mqConfig.Min, mqConfig.Max, func() (mq *MQ, e error) {
		return NewMQ(addr).Open()
	})
	return err
}

//关闭连接池
func CloseMqPool() {
	_ = pool.ClosePool()
}

//获取连接
func GetMQConnect() (*amqp.Channel, error) {
	// 获取连接
	mq, err := pool.Get()
	if err != nil {
		return nil, err
	}
	defer PutMQConnect(mq)
	if st := mq.State(); st != StateOpened {
		return nil, fmt.Errorf("MQ: Not opened, now state is %d", st)
	}
	// 初始化channel
	ch, err := mq.Channel()
	if err != nil {
		return nil, fmt.Errorf("MQ: Create Channel failed, %v", err)
	}
	return ch, nil
}

//释放连接
func PutMQConnect(mq *MQ) {
	if mq.State() == StateOpened {
		_ = pool.Put(mq)
	} else {
		_ = pool.Close(mq)
	}
}
