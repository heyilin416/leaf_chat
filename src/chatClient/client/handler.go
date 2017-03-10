package client

import (
	"common/msg"
	"github.com/name5566/leaf/log"
)

func handleLogin(args []interface{}) {
	recvMsg := args[0].(*msg.L2C_Login)
	if recvMsg.Err != "" {
		Close()
		log.Error("login is error: %v, please input login cmd", recvMsg.Err)
		return
	}

	userData.AccountId = recvMsg.Id
	Start(recvMsg.FrontAddr)

	sendMsg := &msg.C2F_CheckLogin{Token: recvMsg.Token}
	Client.WriteMsg(sendMsg)
}

func handleCheckLogin(args []interface{}) {
	recvMsg := args[0].(*msg.F2C_CheckLogin)
	if recvMsg.Err != "" {
		Close()
		log.Error("check login is error: %v, please input login cmd", recvMsg.Err)
		return
	}

	if recvMsg.UserId != "" {
		userData.UserId = recvMsg.UserId
		userData.UserName = recvMsg.UserName

		log.Release("%v login success", userData.AccountName)
	} else {
		userData.UserName = userData.AccountName

		sendMsg := &msg.C2F_CreateUser{UserName: userData.UserName}
		Client.WriteMsg(sendMsg)
	}
}

func handleCreateUser(args []interface{}) {
	recvMsg := args[0].(*msg.F2C_CreateUser)
	if recvMsg.Err != "" {
		Close()
		log.Error("create user is error: %v, please input login cmd", recvMsg.Err)
		return
	}

	userData.UserId = recvMsg.UserId

	log.Release("%v login and create user success", userData.AccountName)
}

func handleEnterRoom(args []interface{}) {
	recvMsg := args[0].(*msg.F2C_EnterRoom)
	if recvMsg.Err != "" {
		log.Error("enter room is error: %v", recvMsg.Err)
		return
	}

	log.Release("%v user enter %v room success", userData.UserName, enterRoomName)
}