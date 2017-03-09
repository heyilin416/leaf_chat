package cluster

import (
	"github.com/name5566/leaf/cluster"
	"loginServer/login"
)

func Init()  {
	cluster.SetRoute("UpdateClientCount", login.ChanRPC)
	cluster.SetRoute("CheckToken", login.ChanRPC)
}
