package internal

import (
	"common/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/cluster"
	"frontServer/db/mongodb/userDB"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"frontServer/conf"
	"frontServer/user"
	"frontServer/center"
	"time"
)

func init() {
	msg.Processor.SetHandler(&msg.C2F_CheckLogin{}, handleCheckLogin)
	msg.Processor.SetHandler(&msg.C2F_CreateUser{}, handleCreateUser)
	msg.Processor.SetHandler(&msg.C2F_EnterRoom{}, handleEnterRoom)
	msg.Processor.SetHandler(&msg.C2F_LeaveRoom{}, handleLeaveRoom)
	msg.Processor.SetHandler(&msg.C2F_SendMsg{}, handleSendMsg)
}

func onAgentInit(agent gate.Agent) {

}

func onAgentDestroy(agent gate.Agent)() {
	var accountId bson.ObjectId
	if val, ok := agent.UserData().(bson.ObjectId); ok {
		accountId = val
	} else if userData, ok := agent.UserData().(*user.Data); ok {
		accountId = userData.AccountId

		serverRoomMap := userData.GetServerRoomMap()
		for serverName, roomNames := range serverRoomMap {
			cluster.Go(serverName, "LeaveRoom", userData.Id, roomNames, true)
		}

		cluster.Go("world", "AccountOffline", userData.AccountId)
		center.ChanRPC.Go("UserOffline", userData.Id, agent)
	}

	if accountId.Valid() {
		center.ChanRPC.Go("AccountOffline", accountId, agent)
	}
}

func handleCheckLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_CheckLogin)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.F2C_CheckLogin{}
	accountId, err := cluster.Call1("login", "CheckToken", recvMsg.Token, conf.Server.ServerName)
	if err != nil {
		sendMsg.Err = err.Error()
		agent.WriteMsg(sendMsg)
		return
	}

	for {
		ok, err := center.ChanRPC.Call1("AccountOnline", accountId, agent)
		if err != nil {
			sendMsg.Err = err.Error()
			agent.WriteMsg(sendMsg)
			return
		}

		if ok.(bool) {
			break
		} else {
			time.Sleep(time.Second)
		}
	}

	userData, err := userDB.Get(accountId.(bson.ObjectId))
	if err == nil {
		sendMsg.UserId = userData.Id
		sendMsg.UserName = userData.Name

		agent.SetUserData(user.NewData(userData))
		center.ChanRPC.Go("UserOnline", userData.Id, agent)
	} else if err == mgo.ErrNotFound{
		agent.SetUserData(accountId)
	} else {
		sendMsg.Err = err.Error()
	}
	agent.WriteMsg(sendMsg)
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

	userData := &userDB.UserData{Id: bson.NewObjectId(), AccountId: accountId, Name: recvMsg.UserName}
	err := userDB.Create(userData)
	if err == nil {
		sendMsg.UserId = userData.Id

		agent.SetUserData(user.NewData(userData))
		center.ChanRPC.Go("UserOnline", userData.Id, agent)
	} else {
		sendMsg.Err = err.Error()
	}
	agent.WriteMsg(sendMsg)
}

func handleEnterRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_EnterRoom)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.F2C_EnterRoom{RoomName: recvMsg.RoomName}
	userData, ok := agent.UserData().(*user.Data)
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

	var serverName string
	var msgList interface{}
	ret, err := cluster.Call1("world", "GetRoomInfo", recvMsg.RoomName)
	if err == nil {
		serverName = ret.(string)
		msgList, err = cluster.Call1(serverName, "EnterRoom", userData.Id, conf.Server.ServerName, recvMsg.RoomName)
	}

	if err == nil {
		userData.AddRoom(recvMsg.RoomName, serverName)
		sendMsg.MsgList = msgList.([]*msg.ChatMsg)
	} else {
		sendMsg.Err = err.Error()
	}
	agent.WriteMsg(sendMsg)
}

func handleLeaveRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_LeaveRoom)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.F2C_LeaveRoom{RoomName: recvMsg.RoomName}
	userData, ok := agent.UserData().(*user.Data)
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

	serverName := userData.GetRoomServerName(recvMsg.RoomName)
	if serverName == "" {
		sendMsg.Err = "you is not in this room"
		agent.WriteMsg(sendMsg)
		return
	}

	err := cluster.Call0(serverName, "LeaveRoom", userData.Id, []string{recvMsg.RoomName}, false)
	if err == nil {
		userData.RemoveRoom(recvMsg.RoomName)
	} else {
		sendMsg.Err = err.Error()
	}
	agent.WriteMsg(sendMsg)
}

func handleSendMsg(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_SendMsg)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.F2C_SendMsg{}
	userData := agent.UserData().(*user.Data)
	serverName := userData.GetRoomServerName(recvMsg.RoomName)
	if serverName == "" {
		sendMsg.Err = "you have not in this room"
		agent.WriteMsg(sendMsg)
		return
	}

	err := cluster.Call0(serverName, "SendMsg", userData.UserData.Id, recvMsg.RoomName, recvMsg.Msg)
	if err != nil {
		sendMsg.Err = err.Error()
	}
	agent.WriteMsg(sendMsg)
}