package internal

import (
	"loginServer/db/mongodb/tokenDB"
	"github.com/name5566/leaf/cluster"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/chanrpc"
)

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	handleRpc("CheckToken", CheckToken)
}

func CheckToken(args []interface{}) {
	tokenId := args[0].(bson.ObjectId)
	retFunc := args[1].(chanrpc.GetExternalRetFunc)()
	go func() {
		accountId, err := tokenDB.Check(tokenId)
		retFunc(accountId, err)
	}()
}
