package client

import (
	"crypto/md5"
	"fmt"
	"common/msg"
	"chatClient/conf"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

var (
	userData      = UserData{}
)

type UserData struct {
	AccountId   bson.ObjectId
	AccountName string

	UserId   bson.ObjectId
	UserName string
}

func init() {
	skeleton.RegisterCommand("login", "login account: input accountName, passward", login)
	skeleton.RegisterCommand("enterRoom", "enter room: input roomName", enterRoom)
	skeleton.RegisterCommand("leaveRoom", "enter room: input roomName", leaveRoom)
	skeleton.RegisterCommand("sendMsg", "send msg to room: input roomName, msgContent", sendMsg)
}

func login(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 2 {
		err = errors.New("args len is less than 2")
		return
	}

	name := args[0].(string)
	password := args[1].(string)
	userData.AccountName = name

	Start(conf.Client.LoginAddr)

	hash := md5.Sum([]byte(password))
	strMd5 := fmt.Sprintf("%x", hash)
	msg := &msg.C2L_Login{Name: name, Password: strMd5}
	Client.WriteMsg(msg)
	return
}

func enterRoom(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 1 {
		err = errors.New("args len is less than 1")
		return
	}

	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	roomName := args[0].(string)
	msg := &msg.C2F_EnterRoom{RoomName: roomName}
	Client.WriteMsg(msg)
	return
}

func leaveRoom(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 1 {
		err = errors.New("args len is less than 1")
		return
	}

	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	roomName := args[0].(string)
	msg := &msg.C2F_LeaveRoom{RoomName: roomName}
	Client.WriteMsg(msg)
	return
}

func sendMsg(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 2 {
		err = errors.New("args len is less than 2")
		return
	}

	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	roomName := args[0].(string)
	msgContent := args[1].(string)
	msgContent = fmt.Sprintf("%v say: %v", userData.UserName, msgContent)
	msg := &msg.C2F_SendMsg{RoomName: roomName, Msg: []byte(msgContent)}
	Client.WriteMsg(msg)
	return
}

