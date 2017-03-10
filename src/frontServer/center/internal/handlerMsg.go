package internal

import (
	"reflect"
	"common/msg"
	fmsg "frontServer/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/cluster"
	"frontServer/db/mongodb/userC"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"frontServer/conf"
	"github.com/name5566/leaf/log"
)

func handleMsg(m interface{}, h interface{}) {
	fmsg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.C2F_CheckLogin{}, handleCheckLogin)
	handleMsg(&msg.C2F_CreateUser{}, handleCreateUser)
	handleMsg(&msg.C2F_EnterRoom{}, handleEnterRoom)
}

func handleCheckLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_CheckLogin)
	agent := args[1].(gate.Agent)

	var acountId interface{}
	var err error
	skeleton.Go(func() {
		acountId, err = cluster.Call1("login", "CheckToken", recvMsg.Token)
	}, func() {
		sendMsg := &msg.F2C_CheckLogin{}
		if err == nil {
			userData, err := userC.GetUser(acountId.(bson.ObjectId))
			if err == nil {
				sendMsg.UserId = userData.Id
				sendMsg.UserName = userData.Name

				agent.SetUserData(NewUserInfo(userData))
			} else if err == mgo.ErrNotFound{
				agent.SetUserData(acountId)
			} else {
				sendMsg.Err = err.Error()
			}
		} else {
			sendMsg.Err = err.Error()
		}
		agent.WriteMsg(sendMsg)
	})
}

func handleCreateUser(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_CreateUser)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.F2C_CreateUser{}
	accountId, ok := agent.UserData().(bson.ObjectId)
	if !ok {
		sendMsg.Err = "login step is error"
		agent.WriteMsg(sendMsg)
		return
	}

	if recvMsg.UserName == "" {
		sendMsg.Err = "user name is null"
		agent.WriteMsg(sendMsg)
		return
	}

	var err error
	var userData *userC.UserData
	skeleton.Go(func() {
		userData = &userC.UserData{Id: bson.NewObjectId(), AccountId: accountId, Name: recvMsg.UserName}
		err = userC.CreateUser(userData)
	}, func() {
		if err == nil {
			sendMsg.UserId = userData.Id

			agent.SetUserData(NewUserInfo(userData))
		} else {
			sendMsg.Err = err.Error()
		}
		agent.WriteMsg(sendMsg)
	})
}

func handleEnterRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_EnterRoom)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.F2C_EnterRoom{}
	userInfo, ok := agent.UserData().(*UserInfo)
	if !ok {
		sendMsg.Err = "you is not login success"
		agent.WriteMsg(sendMsg)
		return
	}

	if recvMsg.RoomName == "" {
		sendMsg.Err = "room name is null"
		agent.WriteMsg(sendMsg)
		return
	}

	if _, ok := userInfo.roomMap[recvMsg.RoomName]; ok {
		sendMsg.Err = "you have in this room"
		agent.WriteMsg(sendMsg)
		return
	}

	var err error
	var serverName string
	var msgList interface{}
	skeleton.Go(func() {
		var ret interface{}
		ret, err = cluster.Call1("world", "GetRoomInfo", recvMsg.RoomName)
		if err == nil {
			serverName = ret.(string)
			msgList, err = cluster.Call1(serverName, "EnterRoom", userInfo.Id, conf.Server.ServerName, recvMsg.RoomName)
		}
	}, func() {
		if err == nil {
			userInfo.roomMap[recvMsg.RoomName] = serverName
			sendMsg.MsgList = msgList.([]*msg.ChatMsg)
			log.Debug("%v enter %v room, all rooms %v", userInfo.Name, recvMsg.RoomName, userInfo.roomMap)
		} else {
			sendMsg.Err = err.Error()
		}
		agent.WriteMsg(sendMsg)
	})
}