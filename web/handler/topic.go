package handler

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/model"
	"amqp-agent/common/util"
	"amqp-agent/service"
	"amqp-agent/web/form"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type Topic struct{}

func (t *Topic) Publish(ctx *gin.Context) {
	params := ctx.MustGet(constant.SysBindFormKey).(*form.TopicPublishForm)
	if len(params.TagList) == 0 && params.RoutingKey == "" {
		util.Fail(ctx, constant.ErrorParams)
		return
	}
	messageList, err := util.ParseMessage(params.Data)
	if err != nil {
		util.Fail(ctx, constant.ErrorParams)
		return
	}
	c, _ := util.ContextWithSpan(ctx)

	//延时消息入库
	if params.DelayTime > 0 {
		js, _ := json.Marshal(params)
		if _, err := new(model.DelayMessage).SaveOne(params.Source, constant.MessageTypeTopic, params.TopicName, string(js), params.DelayTime); err != nil {
			util.Fail(ctx, err)
			return
		}
		util.Success(ctx, nil)
		return
	}

	data, err := new(service.Producer).PublishTopic(c, params.TopicName, params.RoutingKey, params.TagList, messageList)
	if err != nil {
		util.Fail(ctx, err)
		return
	}
	util.Success(ctx, data)
}
