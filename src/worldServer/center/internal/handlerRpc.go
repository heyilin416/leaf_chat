package internal

import (
	"github.com/name5566/leaf/cluster"
	"github.com/name5566/leaf/log"
	"github.com/pkg/errors"
)

var (
	frontInfoMap = map[string]*FrontInfo{}
	chatInfoMap  = map[string]*ChatInfo{}
	roomInfoMap  = map[string]*RoomInfo{}
)

type FrontInfo struct {
	serverName     string
	clientCount    int
	maxClientCount int
	clientAddr     string
}

type ChatInfo struct {
	serverName  string
	clientCount int
	clusterAddr string
}

type RoomInfo struct {
	serverName string
}

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	cluster.AgentChanRPC = ChanRPC

	skeleton.RegisterChanRPC("NewServerAgent", NewServerAgent)
	skeleton.RegisterChanRPC("CloseServerAgent", CloseServerAgent)

	handleRpc("GetBestFrontInfo", GetBestFrontInfo)
	handleRpc("UpdateFrontInfo", UpdateFrontInfo)
	handleRpc("UpdateChatInfo", UpdateChatInfo)
	handleRpc("GetRoomInfo", GetRoomInfo)
	handleRpc("DestroyRoom", DestroyRoom)
}

func NewServerAgent(args []interface{}) {
	serverName := args[0].(string)
	agent := args[1].(*cluster.Agent)
	if serverName[:5] == "front" {
		results, err := agent.CallN("GetFrontInfo")
		if err == nil {
			clientCount := results[0].(int)
			maxClientCount := results[1].(int)
			clientAddr := results[2].(string)
			frontInfoMap[serverName] = &FrontInfo{serverName: serverName, clientCount: clientCount,
				maxClientCount: maxClientCount, clientAddr: clientAddr}

			if len(chatInfoMap) > 0 {
				serverInfoMap := map[string]string{}
				for chatName, chatInfo := range chatInfoMap {
					serverInfoMap[chatName] = chatInfo.clusterAddr
				}
				agent.Go("AddClusterClient", serverInfoMap)
			}
		} else {
			log.Error("GetFrontInfo is error: %v", err)
		}
	} else if serverName[:4] == "chat" {
		results, err := agent.CallN("GetChatInfo")
		if err == nil {
			clientCount := results[0].(int)
			clusterAddr := results[1].(string)
			chatInfoMap[serverName] = &ChatInfo{serverName: serverName, clientCount: clientCount, clusterAddr: clusterAddr}

			cluster.Broadcast("front", "AddClusterClient", map[string]string{serverName: clusterAddr})
		} else {
			log.Error("GetChatInfo is error: %v", err)
		}
	}
}

func CloseServerAgent(args []interface{}) {
	serverName := args[0].(string)
	if serverName[:5] == "front" {
		_, ok := frontInfoMap[serverName]
		if ok {
			delete(frontInfoMap, serverName)
		}
	} else if serverName[:4] == "chat" {
		_, ok := chatInfoMap[serverName]
		if ok {
			delete(chatInfoMap, serverName)

			cluster.Broadcast("front", "RemoveClusterClient", serverName)
		}
	}
}

func GetBestFrontInfo(args []interface{}) (addr interface{}, err error) {
	var frontInfo *FrontInfo
	minClientCount := 1<<31 - 1
	for _, _frontInfo := range frontInfoMap {
		if _frontInfo.clientCount < minClientCount && _frontInfo.clientCount < _frontInfo.maxClientCount {
			frontInfo = _frontInfo
		}
	}

	if frontInfo == nil {
		err = errors.New("No front server to alloc")
	} else {
		addr = frontInfo.clientAddr
	}
	return
}

func UpdateFrontInfo(args []interface{}) {
	serverName := args[0].(string)
	clientCount := args[1].(int)
	frontInfo, ok := frontInfoMap[serverName]
	if ok {
		frontInfo.clientCount = clientCount
		log.Debug("%v server of client count is %v", serverName, clientCount)
	}
}

func UpdateChatInfo(args []interface{}) {
	serverName := args[0].(string)
	clientCount := args[1].(int)
	chatInfo, ok := chatInfoMap[serverName]
	if ok {
		chatInfo.clientCount = clientCount
		log.Debug("%v server of client count is %v", serverName, clientCount)
	}
}

func GetRoomInfo(args []interface{}) (serverName interface{}, err error) {
	roomName := args[0].(string)
	roomInfo, ok := roomInfoMap[roomName]
	if ok {
		serverName = roomInfo.serverName
	} else {
		var chatInfo *ChatInfo
		minClientCount := 1<<31 - 1
		for _, _chatInfo := range chatInfoMap {
			if _chatInfo.clientCount < minClientCount {
				chatInfo = _chatInfo
			}
		}

		if chatInfo == nil {
			err = errors.New("No chat server to alloc")
		} else {
			serverName = chatInfo.serverName
			roomInfoMap[roomName] = &RoomInfo{serverName: chatInfo.serverName}
		}
	}
	return
}

func DestroyRoom(args []interface{}) {
	roomName := args[0].(string)
	if _, ok := roomInfoMap[roomName]; ok {
		delete(roomInfoMap, roomName)
		log.Debug("%v room is destroy", roomName)
	}
}