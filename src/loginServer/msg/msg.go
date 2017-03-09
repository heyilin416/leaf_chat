package msg

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

func init() {
	Processor.Register(&C2L_Login{})
	Processor.Register(&L2C_Login{})
}
