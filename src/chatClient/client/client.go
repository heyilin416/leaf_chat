package client

import (
	"github.com/name5566/leaf/network"
	"github.com/name5566/leaf/log"
	"reflect"
	"net"
	"time"
)

var (
	Client		*Agent
	processor 	network.Processor
)

func newAgent(conn *network.TCPConn) network.Agent {
	Client = new(Agent)
	Client.conn = conn
	return Client
}

type Agent struct {
	conn       	*network.TCPConn
	userData   	interface{}
}
func (a *Agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Error("read message: %v", err)
			break
		}

		if processor != nil {
			msg, err := processor.Unmarshal(data)
			if err != nil {
				log.Error("unmarshal message error: %v", err)
				break
			}
			err = processor.Route(msg, a)
			if err != nil {
				log.Error("route message error: %v", err)
				break
			}
		}
	}
}

func (a *Agent) OnClose() {}

func (a *Agent) WriteMsg(msg interface{}) {
	if processor != nil {
		data, err := processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *Agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *Agent) Close() {
	a.conn.Close()
}

func (a *Agent) Destroy() {
	a.conn.Destroy()
}

func (a *Agent) UserData() interface{} {
	return a.userData
}

func (a *Agent) SetUserData(data interface{}) {
	a.userData = data
}

func Init(p network.Processor)  {
	processor = p
}

func Close() {
	if Client != nil {
		Client.Destroy()
		Client = nil
	}
}

func Start(addr string)  {
	Close()

	client := new(network.TCPClient)
	client.Addr = addr
	client.NewAgent = newAgent
	client.Start()

	log.Release("start connect to %v", addr)
	for {
		time.Sleep(time.Second)
		if Client != nil {
			break
		}
	}
	log.Release("connect %v success", addr)
}