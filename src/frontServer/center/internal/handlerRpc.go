package internal

import (
	"frontServer/conf"
	"github.com/name5566/leaf/cluster"
	"github.com/name5566/leaf/gate"
)

var (
	clientCount = 0
)

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	skeleton.RegisterChanRPC("NewAgent", NewAgent)
	skeleton.RegisterChanRPC("CloseAgent", CloseAgent)

	handleRpc("GetFrontInfo", GetFrontInfo)
	handleRpc("AddClusterClient", AddClusterClient)
	handleRpc("RemoveClusterClient", RemoveClusterClient)
}

func NewAgent(args []interface{}) {
	clientCount += 1
	cluster.Go("world", "UpdateFrontInfo", conf.Server.ServerName, clientCount)
}

func CloseAgent(args []interface{}) error {
	agent := args[0].(gate.Agent)
	if userInfo, ok := agent.UserData().(*UserInfo); ok {
		serverRoomMap := map[string][]string{}
		for roomName, serverName := range userInfo.roomMap {
			if roomNames, ok := serverRoomMap[serverName]; ok {
				serverRoomMap[serverName] = append(roomNames, roomName)
			} else {
				serverRoomMap[serverName] = []string{roomName}
			}
		}

		for serverName, roomNames := range serverRoomMap {
			cluster.Go(serverName, "LeaveRoom", userInfo.Id, roomNames)
		}
	}

	clientCount -= 1
	cluster.Go("world", "UpdateFrontInfo", conf.Server.ServerName, clientCount)
	return nil
}

func GetFrontInfo(args []interface{}) ([]interface{}, error) {
	return []interface{}{clientCount, conf.Server.MaxConnNum, conf.Server.TCPAddr}, nil
}

func AddClusterClient(args []interface{}) {
	serverInfoMap := args[0].(map[string]string)
	for serverName, addr := range serverInfoMap {
		cluster.AddClient(serverName, addr)
	}
}

func RemoveClusterClient(args []interface{}) {
	serverName := args[0].(string)
	cluster.RemoveClient(serverName)
}