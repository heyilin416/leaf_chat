package internal

import (
	"github.com/name5566/leaf/module"
	"chatServer/base"
	"github.com/name5566/leaf/chanrpc"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/cluster"
	"github.com/name5566/leaf/log"
	"common/msg"
	"time"
	"chatServer/conf"
	"chatServer/common"
	"errors"
	"sync/atomic"
	"fmt"
)

func NewModule(id int) *Module {
	skeleton := base.NewSkeleton()
	module := &Module{Skeleton: skeleton, ChanRPC: skeleton.ChanRPCServer}
	module.name = fmt.Sprintf("module%v", id)
	module.roomInfoMap = map[string]*RoomInfo{}

	RegisterHandler(module.ChanRPC)
	return module
}

type Module struct {
	*module.Skeleton
	ChanRPC *chanrpc.Server

	name          string
	clientCount   int32
	roomInfoMap   map[string]*RoomInfo
}

func (m *Module) OnInit() {
	m.Skeleton.AfterFunc(time.Duration(conf.DestroyRoomInterval/10), m.checkDestroyRoom)
}

func (m *Module) OnDestroy() {

}

func (m *Module) checkDestroyRoom() {
	destroyRooms := []string{}
	nowTime := time.Now().Unix()
	for roomName, roomInfo := range m.roomInfoMap {
		if roomInfo.CheckDestroy(nowTime) {
			destroyRooms = append(destroyRooms, roomName)
		}
	}

	if len(destroyRooms) > 0 {
		for _, roomName := range destroyRooms {
			delete(m.roomInfoMap, roomName)
		}
		cluster.Go("world", "DestroyRoom", destroyRooms)
		log.Debug("%v rooms is destroy in %v", destroyRooms, m.name)
	}

	m.Skeleton.AfterFunc(time.Duration(conf.DestroyRoomInterval/10), m.checkDestroyRoom)
}

func (m *Module) GetChanRPC() *chanrpc.Server {
	return m.ChanRPC
}

func (m *Module) GetClientCount() int32 {
	return atomic.LoadInt32(&m.clientCount)
}

func (m *Module) AddClientCount(delta int32) {
	atomic.AddInt32(&m.clientCount, delta)
	common.AddClientCount(delta)
}

func (m *Module) GetRoomInfo(roomName string) *RoomInfo {
	return m.roomInfoMap[roomName]
}

func (m *Module) NewRoom(name string) *RoomInfo {
	room := &RoomInfo{name: name}
	room.module = m
	room.userServerMap = map[bson.ObjectId]string{}
	m.roomInfoMap[name] = room
	return room
}

type RoomInfo struct {
	module *Module

	name          string
	startNullTime int64
	msgList       []*msg.ChatMsg
	userServerMap map[bson.ObjectId]string
}

func (r *RoomInfo) CheckDestroy(curTime int64) bool {
	if r.GetUserCount() < 1 {
		if r.startNullTime > 0 {
			if curTime-r.startNullTime >= int64(conf.DestroyRoomInterval) {
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
		r.module.AddClientCount(1)
		log.Debug("%v user enter %v room, %v count user in %v", userId, r.name, r.GetUserCount(), r.module.name)
	}
	return ok
}

func (r *RoomInfo) LeaveRoom(userId bson.ObjectId) bool {
	ok := r.IsInRoom(userId)
	if ok {
		delete(r.userServerMap, userId)
		r.module.AddClientCount(-1)
		log.Debug("%v user leave %v room, %v count user in %v", userId, r.name, r.GetUserCount(), r.module.name)
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
