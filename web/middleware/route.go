package middleware

import (
	"amqp-agent/common/constant"
	"amqp-agent/common/util"
	"amqp-agent/web/form"
	"amqp-agent/web/handler"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"reflect"
)

func RegisterRoute(engine *gin.Engine) {
	route := engine.Group(constant.SysRoutePre)
	{
		route.GET("/queue/pull", new(handler.Queue).Pull)
		route.POST("/queue/ack", bind(form.QueueAckForm{}), new(handler.Queue).Ack)
		route.POST("/queue/send", bind(form.QueueSendForm{}), new(handler.Queue).Send)
		route.POST("/topic/publish", bind(form.TopicPublishForm{}), new(handler.Topic).Publish)
	}
}

/**
 * 绑定路由参数
 */
func bind(val interface{}) gin.HandlerFunc {
	value := reflect.ValueOf(val)
	if value.Kind() == reflect.Ptr {
		panic(`Bind bean can not be a pointer. Example: Use: gin.Bind(Struct{}) instead of gin.Bind(&Struct{})`)
	}

	return func(c *gin.Context) {
		obj := reflect.New(value.Type()).Interface()
		// shouldBindBodyWith 在context存储了 BodyBytesKey, 需要在后置的middleware当中使用该值
		err := c.ShouldBindBodyWith(obj, binding.JSON)
		if err != nil {
			util.Fail(c, constant.ErrorParams)
			c.Abort()
			return
		}
		c.Set(constant.SysBindFormKey, obj)
		c.Next()
	}
}
