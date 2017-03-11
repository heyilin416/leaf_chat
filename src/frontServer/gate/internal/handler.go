package internal

import (
	"common/msg"
	fmsg "frontServer/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/cluster"
	"frontServer/db/mongodb/userDB"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"frontServer/conf"
	"github.com/name5566/leaf/log"
	"frontServer/user"
	"frontServer/center"
	"time"
)

func init() {
	fmsg.Processor.SetHandler(&msg.C2F_CheckLogin{}, handleCheckLogin)
	fmsg.Processor.SetHandler(&msg.C2F_CreateUser{}, handleCreateUser)
	fmsg.Processor.SetHandler(&msg.C2F_EnterRoom{}, handleEnterRoom)
}

func EndAgentRun(agent gate.Agent)() {
	var accountId bson.ObjectId
	if val, ok := agent.UserData().(bson.ObjectId); ok {
		accountId = val
	} else if userData, ok := agent.UserData().(*user.Data); ok {
		accountId = userData.AccountId
		serverRoomMap := map[string][]string{}
		for roomName, serverName := range userData.RoomMap {
			if roomNames, ok := serverRoomMap[serverName]; ok {
				serverRoomMap[serverName] = append(roomNames, roomName)
			} else {
				serverRoomMap[serverName] = []string{roomName}
			}
		}

		for serverName, roomNames := range serverRoomMap {
			cluster.Go(serverName, "LeaveRoom", userData.Id, roomNames)
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
	accountId, err := cluster.Call1("login", "CheckToken", recvMsg.Token)
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

	sendMsg := &msg.F2C_EnterRoom{}
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

	if _, ok := userData.RoomMap[recvMsg.RoomName]; ok {
		sendMsg.Err = "you have in this room"
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
		userData.RoomMap[recvMsg.RoomName] = serverName
		sendMsg.MsgList = msgList.([]*msg.ChatMsg)
		log.Debug("%v enter %v room, all rooms %v", userData.Name, recvMsg.RoomName, userData.RoomMap)
	} else {
		sendMsg.Err = err.Error()
	}
	agent.WriteMsg(sendMsg)
}