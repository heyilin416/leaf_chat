package cluster

import (
	"github.com/name5566/leaf/cluster"
	"frontServer/center"
)

func Init()  {
	cluster.SetRoute("GetFrontInfo", center.ChanRPC)
}
