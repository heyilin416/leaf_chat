package user

import (
	"frontServer/db/mongodb/userDB"
	"github.com/name5566/leaf/log"
)

type Data struct {
	*userDB.UserData
	roomMap map[string]string
}

func NewData(userData *userDB.UserData) *Data {
	return &Data{UserData: userData, roomMap: map[string]string{}}
}

func (a *Data) GetRoomServerName(roomName string) string {
	for _roomName, serverName := range a.roomMap {
		if _roomName == roomName {
			return serverName
		}
	}
	return ""
}

func (a *Data) IsInRoom(roomName string) bool {
	return a.GetRoomServerName(roomName) != ""
}

func (a *Data) AddRoom(roomName, serverName string) {
	a.roomMap[roomName] = serverName
	log.Debug("%v enter %v room, all rooms %v", a.Name, roomName, a.roomMap)
}

func (a *Data) RemoveRoom(roomName string) {
	if _, ok := a.roomMap[roomName]; ok {
		delete(a.roomMap, roomName)
	}
}

func (a *Data) GetServerRoomMap() map[string][]string {
	serverRoomMap := map[string][]string{}
	for roomName, serverName := range a.roomMap {
		if roomNames, ok := serverRoomMap[serverName]; ok {
			serverRoomMap[serverName] = append(roomNames, roomName)
		} else {
			serverRoomMap[serverName] = []string{roomName}
		}
	}
	return serverRoomMap
}
