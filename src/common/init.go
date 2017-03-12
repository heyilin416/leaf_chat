package common

import (
	"common/msg"
	"encoding/gob"
	"gopkg.in/mgo.v2/bson"
)

func Init() {
	gob.Register(bson.NewObjectId())
	gob.Register([]bson.ObjectId{})
	gob.Register(map[string]string{})
	gob.Register(msg.ChatMsg{})
	gob.Register([]*msg.ChatMsg{})
}