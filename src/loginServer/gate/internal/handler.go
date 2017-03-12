package internal

import (
	"common/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/cluster"
	"loginServer/db/mongodb/accountDB"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"loginServer/db/mongodb/tokenDB"
)

func init() {
	msg.Processor.SetHandler(&msg.C2L_Login{}, handleLogin)
}

func handleLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_Login)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.L2C_Login{}
	sendErrFunc := func(err string) {
		sendMsg.Err = err
		agent.WriteMsg(sendMsg)
	}

	if recvMsg.Name == "" {
		sendErrFunc("account name is null")
		return
	}

	accountData, err := accountDB.Get(recvMsg.Name)
	if err == mgo.ErrNotFound {
		accountData = &accountDB.Data{Id: bson.NewObjectId(), Name: recvMsg.Name, Password: recvMsg.Password}
		err = accountDB.Create(accountData)
	}

	if err != nil {
		sendErrFunc(err.Error())
		return
	} else if accountData.Password != recvMsg.Password {
		sendErrFunc("password is error")
		return
	}

	results, err := cluster.CallN("world", "GetBestFrontInfo", accountData.Id)
	if err != nil {
		sendErrFunc(err.Error())
		return
	}

	frontName := results[0].(string)
	frontAddr := results[1].(string)
	token, err := tokenDB.Create(accountData.Id, frontName)
	if err != nil {
		sendErrFunc(err.Error())
		return
	}

	sendMsg.Id = accountData.Id
	sendMsg.FrontAddr = frontAddr
	sendMsg.Token = token
	agent.WriteMsg(sendMsg)
}