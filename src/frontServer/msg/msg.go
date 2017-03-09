package msg

import (
	"github.com/name5566/leaf/network/json"
)

var (
	Processor = json.NewProcessor()
)

type C2F_CheckLogin struct {
	Token string
}

type F2C_CheckLogin struct {
	Err string
}

func init() {
	Processor.Register(&C2F_CheckLogin{})
	Processor.Register(&F2C_CheckLogin{})
}
