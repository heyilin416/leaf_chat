package internal

import (
	"frontServer/db/mongodb/userC"
)

type UserInfo struct {
	*userC.UserData
	roomMap map[string]string
}

func NewUserInfo(userData *userC.UserData) *UserInfo {
	return &UserInfo{UserData: userData, roomMap: map[string]string{}}
}