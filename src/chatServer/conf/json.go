package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
)

var Server struct {
	LogLevel        string
	LogPath         string
	WSAddr          string
	CertFile        string
	KeyFile         string
	TCPAddr         string
	MaxConnNum      int
	ConsolePort     int
	ProfilePath     string
	RoomModuleCount int

	MongodbAddr       string
	MongodbSessionNum int

	ServerName      string
	ListenAddr      string
	ConnAddrs       map[string]string
	PendingWriteNum int
}

func Init(fileName string) {
	data, err := ioutil.ReadFile("conf/" + fileName)
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
