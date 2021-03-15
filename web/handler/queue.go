package handler

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/model"
	"amqp-agent/common/util"
	"amqp-agent/service"
	"amqp-agent/web/form"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Queue struct{}

func (q *Queue) Send(ctx *gin.Context) {
	params := ctx.MustGet(constant.SysBindFormKey).(*form.QueueSendForm)
	messageList, err := util.ParseMessage(params.Data)
	if err != nil {
		util.Fail(ctx, constant.ErrorParams)
		return
	}
	c, _ := util.ContextWithSpan(ctx)

	//延时消息入库
	if params.DelayTime > 0 {
		js, _ := json.Marshal(params)
		if _, err := new(model.DelayMessage).SaveOne(params.Source, constant.MessageTypeQueue, params.QueueName, string(js), params.DelayTime); err != nil {
			util.Fail(ctx, err)
			return
		}
		util.Success(ctx, nil)
		return
	}

	data, err := new(service.Producer).EnQueue(c, params.QueueName, messageList)
	if err != nil {
		util.Fail(ctx, err)
		return
	}
	util.Success(ctx, data)
}

func (q *Queue) Pull(ctx *gin.Context) {
	queueName := ctx.Query("queueName")
	messageCount, _ := strconv.Atoi(ctx.Query("messageCount"))
	blockSecond, _ := strconv.Atoi(ctx.Query("blockSecond"))
	if queueName == "" {
		util.Fail(ctx, constant.ErrorParams)
		return
	}
	if messageCount <= 0 {
		messageCount = 1
	}
	if blockSecond <= 0 {
		blockSecond = 1
	}
	c, _ := util.ContextWithSpan(ctx)

	data, err := new(service.Consumer).DeQueue(c, queueName, messageCount, blockSecond)
	if err != nil {
		util.Fail(ctx, err)
		return
	}

	util.Success(ctx, data)
}

func (q *Queue) Ack(ctx *gin.Context) {
	params := ctx.MustGet(constant.SysBindFormKey).(*form.QueueAckForm)
	c, _ := util.ContextWithSpan(ctx)

	if err := new(service.Consumer).Ack(c, params.QueueName, params.Data); err != nil {
		util.Fail(ctx, err)
		return
	}

	util.Success(ctx, nil)
}
