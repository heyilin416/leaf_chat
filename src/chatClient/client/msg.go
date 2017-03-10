package client

import (
	"github.com/name5566/leaf/network/json"
	"common/msg"
)

var (
	Processor = json.NewProcessor()
)

func init() {
	Processor.Register(&msg.C2L_Login{})
	Processor.Register(&msg.C2F_CheckLogin{})
	Processor.Register(&msg.C2F_CreateUser{})
	Processor.Register(&msg.C2F_EnterRoom{})

	Processor.SetHandler(&msg.L2C_Login{}, handleLogin)
	Processor.SetHandler(&msg.F2C_CheckLogin{}, handleCheckLogin)
	Processor.SetHandler(&msg.F2C_CreateUser{}, handleCreateUser)
	Processor.SetHandler(&msg.F2C_EnterRoom{}, handleEnterRoom)
}
