package msg

import (
	"github.com/name5566/leaf/network/json"
	"common/msg"
)

var (
	Processor = json.NewProcessor()
)

func init() {
	Processor.Register(&msg.C2L_Login{})
	Processor.Register(&msg.L2C_Login{})
}
