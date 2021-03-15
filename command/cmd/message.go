package cmd

import (
	"amqp-agent/command/handler"
	"context"
	"github.com/spf13/cobra"
)

//处理延时队列
var DelayMessage = &cobra.Command{
	Use:   "delay-message",
	Short: "处理延时消息",
	Long:  "略",
	Run: func(cmd *cobra.Command, args []string) {
		new(handler.DelayMessage).Handle(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(DelayMessage)
}
