package internal

import (
	"reflect"
	"worldServer/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/cluster"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.C2F_CheckLogin{}, handleCheckLogin)
}

func handleCheckLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2F_CheckLogin)
	agent := args[1].(gate.Agent)

	var id interface{}
	var err error
	skeleton.Go(func() {
		id, err = cluster.Call1("login", "CheckToken", recvMsg.Token)
	}, func() {
		sendMsg := &msg.F2C_CheckLogin{}
		if err == nil {
			agent.SetUserData(id)
		} else {
			sendMsg.Err = err.Error()
		}
		agent.WriteMsg(sendMsg)
	})
}