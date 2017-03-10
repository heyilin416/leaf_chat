package internal

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"crypto/md5"
	"fmt"
	"math/rand"
	"github.com/name5566/leaf/cluster"
)

var (
	tokenMap = map[string]bson.ObjectId{}
)

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	handleRpc("CheckToken", CheckToken)
}

func createToken(id bson.ObjectId) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%x", rand.Uint64())))
	token := fmt.Sprintf("%x", hash)
	tokenMap[token] = id
	return token
}

func CheckToken(args []interface{}) (id interface{}, err error) {
	token := args[0].(string)

	var ok bool
	id, ok = tokenMap[token]
	if !ok {
		err = errors.New("token is not exist")
	}
	return
}
