package main

import (
	"common"
	"common/msg"
	"chatClient/client"
	"github.com/name5566/leaf"
)

func main()  {
	common.Init()
	client.Init(msg.Processor)

	leaf.Run(
		client.Module,
	)
}
