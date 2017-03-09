package mongodb

import (
	"github.com/name5566/leaf/db/mongodb"
	"frontServer/conf"
	"github.com/name5566/leaf/log"
)

var (
	Context *mongodb.DialContext
)

func init()  {
	var err error
	Context, err = mongodb.Dial(conf.Server.MongodbAddr, conf.Server.MongodbSessionNum)
	if err != nil {
		log.Fatal("mongondb init is error(%v)", err)
	}
}
