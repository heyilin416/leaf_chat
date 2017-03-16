package internal

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/chanrpc"
)

func RegisterHandler(chanRPC *chanrpc.Server) {
	chanRPC.Register("EnterRoom", EnterRoom)
	chanRPC.Register("LeaveRoom", LeaveRoom)
	chanRPC.Register("SendMsg", SendMsg)
}

func EnterRoom(args []interface{}) (interface{}, error) {
	module := args[0].(*Module)
	roomName := args[1].(string)
	userId := args[2].(bson.ObjectId)
	serverName := args[3].(string)

	roomInfo := module.GetRoomInfo(roomName)
	if roomInfo == nil {
		roomInfo = module.NewRoom(roomName)
	}

	if roomInfo.EnterRoom(userId, serverName) {
		return roomInfo.GetMsgList(), nil
	}
	return nil, errors.New("you have in this room")
}

func LeaveRoom(args []interface{}) error {
	module := args[0].(*Module)
	roomName := args[1].(string)
	userId := args[2].(bson.ObjectId)

	roomInfo := module.GetRoomInfo(roomName)
	if roomInfo == nil {
		return errors.New("this room is not exist")
	}

	if !roomInfo.LeaveRoom(userId) {
		return errors.New("you is not in this room")
	}
	return nil
}

func SendMsg(args []interface{}) error {
	module := args[0].(*Module)
	roomName := args[1].(string)
	userId := args[2].(bson.ObjectId)
	msgContent := args[3].([]byte)

	roomInfo := module.GetRoomInfo(roomName)
	if roomInfo == nil {
		return errors.New("this room is not exist")
	}
	return roomInfo.SendMsg(userId, msgContent)
}

