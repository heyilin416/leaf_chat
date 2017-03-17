package internal

import (
	"loginServer/db/mongodb/tokenDB"
	"github.com/name5566/leaf/cluster"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/chanrpc"
)

func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

func init() {
	handleRpc("CheckToken", CheckToken, chanrpc.FuncExtRet)
}

func CheckToken(args []interface{}) {
	tokenId := args[0].(bson.ObjectId)
	frontName := args[1].(string)
	retFunc := args[2].(chanrpc.ExtRetFunc)
	go func() {
		accountId, err := tokenDB.Check(tokenId, frontName)
		retFunc(accountId, err)
	}()
}
