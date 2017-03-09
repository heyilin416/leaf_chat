package internal

import (
	"reflect"
	"loginServer/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/cluster"
	"frontServer/db/mongodb/accountC"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.C2L_Login{}, handleLogin)
}

func handleLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_Login)
	agent := args[1].(gate.Agent)

	sendMsg := &msg.L2C_Login{}
	skeleton.Go(func() {
		account, err := accountC.HasAccount(recvMsg.UserName)
		if err != nil {
			sendMsg.Err = err.Error()
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