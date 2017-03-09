package gate

import (
	"loginServer/msg"
	"loginServer/login"
)

func init() {
	msg.Processor.SetRouter(&msg.C2L_Login{}, login.ChanRPC)
}
