package cluster

import (
	"github.com/name5566/leaf/cluster"
	"worldServer/center"
)

func Init()  {
	cluster.SetRoute("GetBestFrontInfo", center.ChanRPC)
}
