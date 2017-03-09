package client

import (
	"github.com/name5566/leaf/log"
)

func handleLogin(args []interface{}) {
	recvMsg := args[0].(*L2C_Login)
	if recvMsg.Err != "" {
		Close()
		log.Error("login is error: %v, please input login cmd", recvMsg.Err)
		return
	}

	accountInfo.Id = recvMsg.Id
	Start(recvMsg.FrontAddr)

	sendMsg := &C2F_CheckLogin{Token: recvMsg.Token}
	Client.WriteMsg(sendMsg)
}

func handleCheckLogin(args []interface{}) {
	recvMsg := args[0].(*F2C_CheckLogin)
	if recvMsg.Err != "" {
		Close()
		log.Error("check login is error: %v, please input login cmd", recvMsg.Err)
		return
	}

	log.Release("%v login success", accountInfo.UserName)
}