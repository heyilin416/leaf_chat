package client

import (
	"github.com/name5566/leaf/module"
	"loginServer/base"
)

var (
	skeleton = base.NewSkeleton()
	Module   = new(ClientModule)
	ChanRPC  = skeleton.ChanRPCServer
)

type ClientModule struct {
	*module.Skeleton
}

func (m *ClientModule) OnInit() {
	m.Skeleton = skeleton
}

func (m *ClientModule) OnDestroy() {

}
