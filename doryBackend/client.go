package doryBackend

import (
	"fmt"

	golcrypt "github.com/abhishekkr/gol/golcrypt"
	gollog "github.com/abhishekkr/gol/gollog"
)

func HandleClientAuth(goldory golcrypt.Dory, axn string) {
	var err error
	switch axn {
	case "ping":
		gollog.Debug("ping")
		err = goldory.Ping()
		fmt.Println(string(goldory.Value))
	case "set":
		gollog.Debug("set")
		err = goldory.Set()
		fmt.Println("while accessing this key, use token:", goldory.Token)
	case "get":
		gollog.Debug("get")
		err = goldory.Get()
		fmt.Println("response:", string(goldory.Value))
	case "del":
		gollog.Debug("del")
		err = goldory.Del()
		fmt.Println("response:", string(goldory.Value))
	case "list":
		gollog.Debug("list")
		err = goldory.List()
		fmt.Println("response:", string(goldory.Value))
	case "purge":
		gollog.Debug("purge")
		err = goldory.PurgeAll()
		fmt.Println("response:", string(goldory.Value))
	default:
		gollog.Err(fmt.Sprintf("unsupported client action", axn))
	}
	if err != nil {
		gollog.Err(err.Error())
	}
}
