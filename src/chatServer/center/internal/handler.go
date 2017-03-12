package internal

import (
	"common/msg"
	"chatServer/conf"
	"github.com/name5566/leaf/cluster"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/log"
	"errors"
	"time"
)

var (
	clientCount   = 0
	roomInfoMap   = map[string]*RoomInfo{}
	userServerMap = map[bson.ObjectId]string{}
)

func init() {
	skeleton.AfterFunc(time.Duration(conf.DestroyRoomInterval / 10), checkDestroyRoom)
}

func checkDestroyRoom() {
	destroyRooms := []string{}
	nowTime := time.Now().Unix()
	for roomName, roomInfo := range roomInfoMap {
		if roomInfo.CheckDestroy(nowTime) {
			destroyRooms = append(destroyRooms, roomName)
		}
	}

	if len(destroyRooms) > 0 {
		for _, roomName := range destroyRooms {
			delete(roomInfoMap, roomName)
		}
		cluster.Go("world", "DestroyRoom", destroyRooms)
		log.Debug("%v rooms is destroy", destroyRooms)
	}

	skeleton.AfterFunc(time.Duration(conf.DestroyRoomInterval / 10), checkDestroyRoom)
}

type RoomInfo struct {
	name          string
	startNullTime int64
	msgList       []*msg.ChatMsg
	userServerMap map[bson.ObjectId]string
}

func NewRoom(name string) *RoomInfo {
	room := &RoomInfo{name: name}
	room.userServerMap = map[bson.ObjectId]string{}
	return room
}

func (r *RoomInfo) CheckDestroy(curTime int64) bool {
	if r.GetUserCount() < 1 {
		if r.startNullTime > 0 {
			if curTime - r.startNullTime >= int64(conf.DestroyRoomInterval) {
				return true
			}
		} else {
			r.startNullTime = time.Now().Unix()
		}
	} else {
		r.startNullTime = 0
	}
	return false
}

func (r *RoomInfo) GetUserCount() int {
	return len(r.userServerMap)
}

func (r *RoomInfo) IsInRoom(userId bson.ObjectId) bool {
	_, ok := r.userServerMap[userId]
	return ok
}

func (r *RoomInfo) EnterRoom(userId bson.ObjectId, serverName string) bool {
	ok := !r.IsInRoom(userId)
	if ok {
		r.userServerMap[userId] = serverName

		clientCount += 1
		cluster.Go("world", "UpdateChatInfo", conf.Server.ServerName, clientCount)

		log.Debug("%v user enter %v room, %v count user", userId, r.name, r.GetUserCount())
	}
	return ok
}

func (r *RoomInfo) LeaveRoom(userId bson.ObjectId) bool {
	ok := r.IsInRoom(userId)
	if ok {
		delete(r.userServerMap, userId)

		clientCount -= 1
		cluster.Go("world", "UpdateChatInfo", conf.Server.ServerName, clientCount)

		log.Debug("%v user leave %v room, %v count user", userId, r.name, r.GetUserCount())
	}
	return ok
}

func (r *RoomInfo) GetMsgList() []*msg.ChatMsg {
	return r.msgList
}

func (r *RoomInfo) SendMsg(userId bson.ObjectId, msgContent []byte) error {
	if !r.IsInRoom(userId) {
		return errors.New("you have not in room")
	}

	msg := &msg.ChatMsg{RoomName: r.name, UserId: userId, MsgTime: time.Now().Unix(), MsgContent: msgContent}
	r.msgList = append(r.msgList, msg)
	msgLen := len(r.msgList)
	if msgLen > conf.MaxRoomMsgLen {
		r.msgList = r.msgList[msgLen-conf.MaxRoomMsgLen:]
	}

	serverUserMap := map[string][]bson.ObjectId{}
	for userId, serverName := range r.userServerMap {
		if userIds, ok := serverUserMap[serverName]; ok {
			serverUserMap[serverName] = append(userIds, userId)
		} else {
			serverUserMap[serverName] = []bson.ObjectId{userId}
		}
	}
	for serverName, userIds := range serverUserMap {
		cluster.Go(serverName, "BroadcastChatMsg", userIds, msg)
	}
	return nil
}

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	handleRpc("GetChatInfo", GetChatInfo)
	handleRpc("EnterRoom", EnterRoom)
	handleRpc("LeaveRoom", LeaveRoom)
	handleRpc("SendMsg", SendMsg)
}

func GetChatInfo(args []interface{}) ([]interface{}, error) {
	return []interface{}{clientCount, conf.Server.ListenAddr}, nil
}

func EnterRoom(args []interface{}) (interface{}, error) {
	userId := args[0].(bson.ObjectId)
	serverName := args[1].(string)
	roomName := args[2].(string)

	userServerMap[userId] = serverName
	roomInfo, ok := roomInfoMap[roomName]
	if !ok {
		roomInfo = NewRoom(roomName)
		roomInfoMap[roomName] = roomInfo
	}

	if roomInfo.EnterRoom(userId, serverName) {
		return roomInfo.GetMsgList(), nil
	}
	return nil, errors.New("you have in this room")
}

func LeaveRoom(args []interface{}) error {
	userId := args[0].(bson.ObjectId)
	roomNames := args[1].([]string)
	isClose := args[2].(bool)

	if isClose {
		if _, ok := userServerMap[userId]; ok {
			delete(userServerMap, userId)
		}
	}

	for _, roomName := range roomNames {
		roomInfo, ok := roomInfoMap[roomName]
		if !ok {
			return errors.New("this room is not exist")
		}

		if !roomInfo.LeaveRoom(userId) {
			return errors.New("you is not in this room")
		}
	}
	return nil
}

func SendMsg(args []interface{}) error {
	userId := args[0].(bson.ObjectId)
	roomName := args[1].(string)
	msgContent := args[2].([]byte)

	roomInfo, ok := roomInfoMap[roomName]
	if !ok {
		return errors.New("room is not exist")
	}

	return roomInfo.SendMsg(userId, msgContent)
}
