package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
)

var Client struct {
	LoginAddr string
}

func init() {
	data, err := ioutil.ReadFile("conf/client.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Client)
	if err != nil {
		log.Fatal("%v", err)
	}
}
