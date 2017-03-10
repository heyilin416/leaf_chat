package internal

import "fmt"

func init() {
	skeleton.RegisterCommand("getFrontInfo", "return all front info", getFrontInfo)
	skeleton.RegisterCommand("getChatInfo", "return all chat info", getChatInfo)
}

func getFrontInfo(args []interface{}) (ret interface{}, err error) {
	ret = fmt.Sprintf("%s", frontInfoMap)
	return
}

func getChatInfo(args []interface{}) (ret interface{}, err error) {
	ret = fmt.Sprintf("%s", chatInfoMap)
	return
}