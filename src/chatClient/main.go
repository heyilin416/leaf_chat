package main

import (
	"chatClient/client"
	"github.com/name5566/leaf"
)

func main()  {
	client.Init(client.Processor)

	leaf.Run(
		client.Module,
	)
}
