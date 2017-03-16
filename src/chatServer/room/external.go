package room

import (
	"chatServer/room/internal"
	"chatServer/conf"
	"github.com/name5566/leaf/module"
	"math"
	"github.com/name5566/leaf/chanrpc"
)

var (
	modules = []*internal.Module{}
)

type Module interface {
	GetChanRPC() *chanrpc.Server
}

func CreateModules() []module.Module {
	results := []module.Module{}
	for i := 0; i < conf.Server.RoomModuleCount; i++ {
		module := internal.NewModule(i)
		modules = append(modules, module)
		results = append(results, module)
	}
	return results
}

func GetBestModule() (module *internal.Module) {
	var minCount int32 = math.MaxInt32
	for _, _module := range modules {
		count := _module.GetClientCount()
		if count < minCount {
			module = _module
			minCount = count
		}
	}
	return
}