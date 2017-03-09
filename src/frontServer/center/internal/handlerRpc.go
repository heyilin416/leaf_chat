package internal

import (
	"frontServer/conf"
	"github.com/name5566/leaf/cluster"
	"gopkg.in/mgo.v2/bson"
	"encoding/gob"
)

var (
	clientCount = 0
)

func init() {
	gob.Register(bson.NewObjectId())

	skeleton.RegisterChanRPC("NewAgent", NewAgent)
	skeleton.RegisterChanRPC("CloseAgent", CloseAgent)
	skeleton.RegisterChanRPC("GetFrontInfo", GetFrontInfo)
}

func NewAgent(args []interface{}) {
	clientCount += 1
	cluster.Go("world", "UpdateFrontInfo", conf.Server.ServerName, clientCount)
}

func CloseAgent(args []interface{}) error {
	clientCount -= 1
	cluster.Go("world", "UpdateFrontInfo", conf.Server.ServerName, clientCount)
	return nil
}

func GetFrontInfo(args []interface{}) ([]interface{}, error) {
	return []interface{}{clientCount, conf.Server.MaxConnNum, conf.Server.TCPAddr}, nil
}
