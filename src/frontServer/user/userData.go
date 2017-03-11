package user

import (
	"frontServer/db/mongodb/userDB"
)

type Data struct {
	*userDB.UserData
	RoomMap map[string]string
}

func NewData(userData *userDB.UserData) *Data {
	return &Data{UserData: userData, RoomMap: map[string]string{}}
}