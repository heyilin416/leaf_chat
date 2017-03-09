package client

import (
	"crypto/md5"
	"fmt"
	"chatClient/conf"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

var (
	accountInfo = AccountInfo{}
)

type AccountInfo struct {
	Id       bson.ObjectId
	UserName string
}

func init() {
	skeleton.RegisterCommand("login", "login account:username passward", login)
}

func login(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 2 {
		err = errors.New("args len is less than 2")
		return
	}

	userName := args[0].(string)
	password := args[1].(string)
	accountInfo.UserName = userName

	Start(conf.Client.LoginAddr)

	hash := md5.Sum([]byte(password))
	strMd5 := fmt.Sprintf("%x", hash)
	msg := &C2L_Login{UserName: userName, Password: strMd5}
	Client.WriteMsg(msg)
	return
}
