package msg

import (
	"gopkg.in/mgo.v2/bson"
)

type ChatMsg struct {
	RoomName   string
	MsgTime    int64
	MsgContent string
}

type C2L_Login struct {
	Name     string
	Password string
}

type L2C_Login struct {
	Err       string
	Id        bson.ObjectId
	FrontAddr string
	Token     string
}

type C2F_CheckLogin struct {
	Token string
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
	Err     string
	MsgList []*ChatMsg
}

type C2F_LeaveRoom struct {
	RoomName string
}

type F2C_LeaveRoom struct {
	Err string
}

type C2F_SendMsg struct {
	Msg string
}

type F2C_SendMsg struct {
	Err string
}

type F2C_MsgList struct {
	MsgList []*ChatMsg
}
