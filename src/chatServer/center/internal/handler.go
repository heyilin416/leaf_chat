package internal

import (
	"chatServer/conf"
	"chatServer/common"
	"chatServer/room"
	"github.com/name5566/leaf/cluster"
	"errors"
	"github.com/name5566/leaf/chanrpc"
)

var (
	roomModuleMap = map[string]room.Module{}
)

func handleRpc(id interface{}, f interface{}, isExtRet bool) {
	cluster.SetRoute(id, ChanRPC)
	if isExtRet {
		ChanRPC.RegisterExtRet(id, f)
	} else {
		ChanRPC.Register(id, f)
	}
}

func init() {
	handleRpc("GetChatInfo", GetChatInfo, false)
	handleRpc("EnterRoom", EnterRoom, true)
	handleRpc("LeaveRoom", LeaveRoom, true)
	handleRpc("LeaveRooms", LeaveRooms, false)
	handleRpc("SendMsg", SendMsg, true)
}

func GetChatInfo(args []interface{}) ([]interface{}, error) {
	return []interface{}{common.GetClientCount(), conf.Server.ListenAddr}, nil
}

func EnterRoom(args []interface{}) {
	roomName := args[0].(string)

	module := roomModuleMap[roomName]
	if module == nil {
		module = room.GetBestModule()
		if module == nil {
			retFunc := args[len(args)-1].(chanrpc.ExtRetFunc)
			retFunc(nil, errors.New("get best room module rpc is fail"))
			return
		}

		roomModuleMap[roomName] = module
	}

	newArgs := []interface{}{module}
	newArgs = append(newArgs, args...)
	skeleton.AsynCall(module.GetChanRPC(), "EnterRoom", newArgs...)
}

func LeaveRoom(args []interface{}) {
	roomName := args[0].(string)

	module := roomModuleMap[roomName]
	if module == nil {
		retFunc := args[len(args)-1].(chanrpc.ExtRetFunc)
		retFunc(nil, errors.New("this room is not exist"))
		return
	}

	newArgs := []interface{}{module}
	newArgs = append(newArgs, args...)
	skeleton.AsynCall(module.GetChanRPC(), "LeaveRoom", newArgs...)
}

func LeaveRooms(args []interface{}) {
	roomNames := args[0].([]string)

	for _, roomName := range roomNames {
		module := roomModuleMap[roomName]
		if module != nil {
			newArgs := []interface{}{module, roomName}
			newArgs = append(newArgs, args[1:]...)
			module.GetChanRPC().Go("LeaveRoom", newArgs...)
		}
	}
}

func SendMsg(args []interface{}) {
	roomName := args[0].(string)

	module := roomModuleMap[roomName]
	if module == nil {
		retFunc := args[len(args)-1].(chanrpc.ExtRetFunc)
		retFunc(nil, errors.New("this room is not exist"))
		return
	}

	newArgs := []interface{}{module}
	newArgs = append(newArgs, args...)
	skeleton.AsynCall(module.GetChanRPC(), "SendMsg", newArgs...)
}
