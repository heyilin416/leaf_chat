package internal

import (
	"reflect"
	"common/msg"
	lmsg "loginServer/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/cluster"
	"loginServer/db/mongodb/accountC"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func handleMsg(m interface{}, h interface{}) {
	lmsg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.C2L_Login{}, handleLogin)
}

func handleLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_Login)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.L2C_Login{}
	if recvMsg.Name == "" {
		sendMsg.Err = "account name is null"
		agent.WriteMsg(sendMsg)
		return
	}

	skeleton.Go(func() {
		account, err := accountC.GetAccount(recvMsg.Name)
		if err == mgo.ErrNotFound {
			account = &accountC.AccountData{Id: bson.NewObjectId(), Name: recvMsg.Name, Password: recvMsg.Password}
			err = accountC.CreateAccount(account)
		}

		if err != nil {
			sendMsg.Err = err.Error()
			return
		} else if account.Password != recvMsg.Password {
			sendMsg.Err = "password is error"
			return
		}

		frontAddr, err := cluster.Call1("world", "GetBestFrontInfo")
		if err != nil {
			sendMsg.Err = err.Error()
			return
		}

		sendMsg.Id = account.Id
		sendMsg.FrontAddr = frontAddr.(string)
	}, func() {
		if sendMsg.Err == "" {
			sendMsg.Token = createToken(sendMsg.Id)
		}
		agent.WriteMsg(sendMsg)
	})
}