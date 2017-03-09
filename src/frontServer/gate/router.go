package gate

import (
	"frontServer/msg"
	"frontServer/center"
)

func init() {
	msg.Processor.SetRouter(&msg.C2F_CheckLogin{}, center.ChanRPC)
}
