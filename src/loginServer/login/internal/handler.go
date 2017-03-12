package internal

import (
	"loginServer/db/mongodb/tokenDB"
	"github.com/name5566/leaf/cluster"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/chanrpc"
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
	handleRpc("CheckToken", CheckToken, true)
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
