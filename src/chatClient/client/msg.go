package client

import (
	"github.com/name5566/leaf/network/json"
	"gopkg.in/mgo.v2/bson"
)

var (
	Processor = json.NewProcessor()
)

type C2L_Login struct {
	UserName string
	Password string
}

type L2C_Login struct {
	Id        bson.ObjectId
	FrontAddr string
	Token     string
	Err       string
}

type C2F_CheckLogin struct {
	Token string
}

type F2C_CheckLogin struct {
	Err string
}

func init() {
	Processor.Register(&C2L_Login{})
	Processor.Register(&C2F_CheckLogin{})

	Processor.SetHandler(&L2C_Login{}, handleLogin)
	Processor.SetHandler(&F2C_CheckLogin{}, handleCheckLogin)
}
