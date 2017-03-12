package client

import (
	"common/msg"
	"github.com/name5566/leaf/log"
	"time"
)

func init() {
	msg.Processor.SetHandler(&msg.L2C_Login{}, handleLogin)
	msg.Processor.SetHandler(&msg.F2C_CheckLogin{}, handleCheckLogin)
	msg.Processor.SetHandler(&msg.F2C_CreateUser{}, handleCreateUser)
	msg.Processor.SetHandler(&msg.F2C_EnterRoom{}, handleEnterRoom)
	msg.Processor.SetHandler(&msg.F2C_LeaveRoom{}, handleLeaveRoom)
	msg.Processor.SetHandler(&msg.F2C_SendMsg{}, handleSendMsg)
	msg.Processor.SetHandler(&msg.F2C_MsgList{}, handleMsgList)
}

func showMsgList(msgList []*msg.ChatMsg) {
	for _, msg := range msgList {
		strTime := time.Unix(msg.MsgTime, 0).Format("2006-01-02 15:04:05")
		log.Release("%v : %v room: %v", strTime, msg.RoomName, string(msg.MsgContent))
	}
}

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

		log.Release("%v(%v) login and create user success", userData.UserName, userData.UserId)
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

	log.Release("%v(%v) login and create user success", userData.UserName, userData.UserId)
}

func handleEnterRoom(args []interface{}) {
	recvMsg := args[0].(*msg.F2C_EnterRoom)
	if recvMsg.Err != "" {
		log.Error("enter %v room is error: %v", recvMsg.RoomName, recvMsg.Err)
		return
	}

	log.Release("you success enter %v room", recvMsg.RoomName)

	showMsgList(recvMsg.MsgList)
}

func handleLeaveRoom(args []interface{}) {
	recvMsg := args[0].(*msg.F2C_LeaveRoom)
	if recvMsg.Err != "" {
		log.Error("leave %v room is error: %v", recvMsg.RoomName, recvMsg.Err)
		return
	}

	log.Release("you success leave %v room", recvMsg.RoomName)
}

func handleSendMsg(args []interface{}) {
	recvMsg := args[0].(*msg.F2C_SendMsg)
	if recvMsg.Err != "" {
		log.Error("send msg is error: %v", recvMsg.Err)
		return
	}

	log.Release("send msg is success")
}

func handleMsgList(args []interface{}) {
	recvMsg := args[0].(*msg.F2C_MsgList)
	showMsgList(recvMsg.MsgList)
}