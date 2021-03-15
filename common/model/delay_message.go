package model

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/util"
)

type DelayMessage struct {
	Base
	DelayMessageId       int64  `gorm:"delay_message_id;primary_key" json:"delay_message_id"`
	DelayMessageSource   string `gorm:"delay_message_source" json:"delay_message_source"`                               // 消息来源
	DelayMessageType     string `gorm:"delay_message_type" json:"delay_message_type"`                                   // 消息类型 QUEUE：队列 TOPIC：主题
	DelayMessageTarget   string `gorm:"delay_message_target" json:"delay_message_target"`                               // 消息目标队列/主题
	DelayMessageStatus   string `gorm:"delay_message_status" json:"delay_message_status"`                               // 消息状态 INIT：初始状态 SUCCESS：成功 FAIL：失败
	DelayMessageData     string `gorm:"delay_message_data" json:"delay_message_data"`                                   // 消息数据
	DelayMessageExtend   string `gorm:"delay_message_extend" json:"delay_message_extend"`                               // 扩展数据
	DelayMessageDelayAt  string `gorm:"delay_message_delay_at;default:CURRENT_TIMESTAMP" json:"delay_message_delay_at"` // 延时发送时间
	DelayMessageCreateAt string `gorm:"delay_message_create_at;default:CURRENT_TIMESTAMP" json:"delay_message_create_at"`
	DelayMessageUpdateAt string `gorm:"delay_message_update_at;default:CURRENT_TIMESTAMP" json:"delay_message_update_at"`
}

func (d *DelayMessage) TableName() string {
	return "delay_message"
}

//保存记录
func (d *DelayMessage) SaveOne(messageSource, messageType, messageTarget, messageData string, delayTime int) (*DelayMessage, error) {
	delayAt := util.GetDateTime(int64(delayTime), constant.DateTimeLayout)
	d.DelayMessageSource = messageSource
	d.DelayMessageType = messageType
	d.DelayMessageTarget = messageTarget
	d.DelayMessageData = messageData
	d.DelayMessageDelayAt = delayAt
	d.DelayMessageStatus = constant.MessageStatusInit
	d.DelayMessageExtend = ""
	if err := d.GetDB().Save(d).Error; err != nil {
		return nil, err
	}

	return d, nil
}

//更新记录
func (d *DelayMessage) UpdateStatus(status, extend string) error {
	return d.GetDB().Model(d).Update(DelayMessage{
		DelayMessageStatus: status,
		DelayMessageExtend: extend,
	}).Error
}

//获取列表
func (d *DelayMessage) GetList() ([]*DelayMessage, error) {
	begin, end := util.GetDateTime(-1*constant.TimeDay, constant.DateTimeLayout), util.Now()
	var delayList []*DelayMessage
	if err := d.GetDB().Where("(delay_message_delay_at between ? and ?) and delay_message_status = ?", begin, end, constant.MessageStatusInit).Order("delay_message_delay_at asc").Limit(5000).Find(&delayList).Error; err != nil {
		return nil, err
	}
	if len(delayList) == 0 {
		return nil, nil
	}
	return delayList, nil
}
