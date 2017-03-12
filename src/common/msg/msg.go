package msg

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/network/json"
)

var (
	Processor = json.NewProcessor()
)

func init() {
	Processor.Register(&C2L_Login{})
	Processor.Register(&L2C_Login{})
	Processor.Register(&C2F_CheckLogin{})
	Processor.Register(&F2C_CheckLogin{})
	Processor.Register(&C2F_CreateUser{})
	Processor.Register(&F2C_CreateUser{})
	Processor.Register(&C2F_EnterRoom{})
	Processor.Register(&F2C_EnterRoom{})
	Processor.Register(&C2F_LeaveRoom{})
	Processor.Register(&F2C_LeaveRoom{})
	Processor.Register(&C2F_SendMsg{})
	Processor.Register(&F2C_SendMsg{})
	Processor.Register(&F2C_MsgList{})
}

type ChatMsg struct {
	RoomName   string
	UserId 	   bson.ObjectId
	MsgTime    int64
	MsgContent []byte
}

type C2L_Login struct {
	Name     string
	Password string
}

type L2C_Login struct {
	Err       string
	Id        bson.ObjectId
	FrontAddr string
	Token     bson.ObjectId
}

type C2F_CheckLogin struct {
	Token bson.ObjectId
}

type F2C_CheckLogin struct {
	Err      string
	UserId   bson.ObjectId
	UserName string
}

type C2F_CreateUser struct {
	UserName string
}

type F2C_CreateUser struct {
	Err    string
	UserId bson.ObjectId
}

type C2F_EnterRoom struct {
	RoomName string
}

type F2C_EnterRoom struct {
	Err      string
	RoomName string
	MsgList  []*ChatMsg
}

type C2F_LeaveRoom struct {
	RoomName string
}

type F2C_LeaveRoom struct {
	Err      string
	RoomName string
}

type C2F_SendMsg struct {
	RoomName string
	Msg      []byte
}

type F2C_SendMsg struct {
	Err string
}

type F2C_MsgList struct {
	MsgList []*ChatMsg
}
