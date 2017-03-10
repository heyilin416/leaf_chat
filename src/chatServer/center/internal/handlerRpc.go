package internal

import (
	"common/msg"
	"chatServer/conf"
	"github.com/name5566/leaf/cluster"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/log"
)

var (
	clientCount   = 0
	roomInfoMap   = map[string]*RoomInfo{}
	userServerMap = map[bson.ObjectId]string{}
)

type RoomInfo struct {
	Name    string
	MsgList []*msg.ChatMsg
	UserIds map[bson.ObjectId]bool
}

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	handleRpc("GetChatInfo", GetChatInfo)
	handleRpc("EnterRoom", EnterRoom)
	handleRpc("LeaveRoom", LeaveRoom)
}

func GetChatInfo(args []interface{}) ([]interface{}, error) {
	return []interface{}{clientCount, conf.Server.ListenAddr}, nil
}

func EnterRoom(args []interface{}) (interface{}, error) {
	userId := args[0].(bson.ObjectId)
	serverName := args[1].(string)
	roomName := args[2].(string)
	roomInfo, ok := roomInfoMap[roomName]
	if !ok {
		userServerMap[userId] = serverName
		roomInfo = &RoomInfo{Name: roomName, UserIds: map[bson.ObjectId]bool{userId: true}}
		roomInfoMap[roomName] = roomInfo
		log.Debug("%v user from %v server enter %v room success, has %v user", userId, serverName, roomName, len(roomInfo.UserIds))

		clientCount += 1
		cluster.Go("world", "UpdateChatInfo", conf.Server.ServerName, clientCount)
	}
	return roomInfo.MsgList, nil
}

func LeaveRoom(args []interface{}) {
	userId := args[0].(bson.ObjectId)
	roomNames := args[1].([]string)

	if _, ok := userServerMap[userId]; ok {
		delete(userServerMap, userId)
	}

	for _, roomName := range roomNames {
		roomInfo, ok := roomInfoMap[roomName]
		if ok {
			if _, ok := roomInfo.UserIds[userId]; ok {
				delete(roomInfo.UserIds, userId)
				log.Debug("%v user leave %v room", userId, roomName)

				clientCount -= 1
				cluster.Go("world", "UpdateChatInfo", conf.Server.ServerName, clientCount)

				if len(roomInfo.UserIds) < 1 {
					delete(roomInfoMap, roomName)
					cluster.Go("world", "DestroyRoom", roomName)
					log.Debug("%v room is destroy", roomName)
				}
			}
		}
	}
}
