package main

import (
	"common"
	"chatClient/client"
	"github.com/name5566/leaf"
)

func main()  {
	common.Init()
	client.Init(client.Processor)

	leaf.Run(
		client.Module,
	)
}
