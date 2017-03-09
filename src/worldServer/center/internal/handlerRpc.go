package internal

import (
	"github.com/name5566/leaf/cluster"
	"github.com/name5566/leaf/log"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"encoding/gob"
)

var (
	frontAgentMap = map[string]*FrontInfo{}
)

type FrontInfo struct {
	clientCount    int
	maxClientCount int
	frontAddr      string
}

func init() {
	cluster.AgentChanRPC = ChanRPC

	gob.Register(bson.NewObjectId())

	skeleton.RegisterChanRPC("NewServerAgent", NewServerAgent)
	skeleton.RegisterChanRPC("CloseServerAgent", CloseServerAgent)
	skeleton.RegisterChanRPC("GetBestFrontInfo", GetBestFrontInfo)
	skeleton.RegisterChanRPC("UpdateFrontInfo", UpdateFrontInfo)
}

func NewServerAgent(args []interface{}) {
	serverName := args[0].(string)
	agent := args[1].(*cluster.Agent)
	if serverName[:5] == "front" {
		results, err := agent.CallN("GetFrontInfo")
		if err == nil {
			clientCount := results[0].(int)
			maxClientCount := results[1].(int)
			frontAddr := results[2].(string)
			frontAgentMap[serverName] = &FrontInfo{clientCount: clientCount, maxClientCount: maxClientCount, frontAddr: frontAddr}
		} else {
			log.Error("GetFrontInfo is error: %v", err)
		}
	}
}

func CloseServerAgent(args []interface{}) {
	serverName := args[0].(string)
	if serverName[:5] == "front" {
		_, ok := frontAgentMap[serverName]
		if ok {
			delete(frontAgentMap, serverName)
		}
	}
}

func GetBestFrontInfo(args []interface{}) (addr interface{}, err error) {
	var frontInfo *FrontInfo
	minClientCount := 1<<31 - 1
	for _, _frontInfo := range frontAgentMap {
		if _frontInfo.clientCount < minClientCount {
			frontInfo = _frontInfo
		}
	}

	if frontInfo == nil {
		err = errors.New("No front server to alloc")
	} else {
		addr = frontInfo.frontAddr
		frontInfo.clientCount += 1
	}
	return
}

func UpdateFrontInfo(args []interface{}) {
	serverName := args[0].(string)
	clientCount := args[1].(int)
	frontInfo, ok := frontAgentMap[serverName]
	if ok {
		frontInfo.clientCount = clientCount
	}
}
