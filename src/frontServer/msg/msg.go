package msg

import (
	"github.com/name5566/leaf/network/json"
	"common/msg"
)

var (
	Processor = json.NewProcessor()
)

func init() {
	Processor.Register(&msg.C2F_CheckLogin{})
	Processor.Register(&msg.F2C_CheckLogin{})
	Processor.Register(&msg.C2F_CreateUser{})
	Processor.Register(&msg.F2C_CreateUser{})
	Processor.Register(&msg.C2F_EnterRoom{})
	Processor.Register(&msg.F2C_EnterRoom{})
}
