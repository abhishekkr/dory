package doryClient

import (
	"fmt"
	"log"
)

func HandleClientAuth(goldory DoryClient, axn string) {
	var err error
	switch axn {
	case "ping":
		err = goldory.Ping()
		fmt.Println(string(goldory.Value))
	case "set":
		err = goldory.Set()
		fmt.Println("while accessing this key, use token:", goldory.Token)
	case "get":
		err = goldory.Get()
		fmt.Println(string(goldory.Value))
	case "del":
		err = goldory.Del()
		fmt.Println(string(goldory.Value))
	case "list":
		err = goldory.List()
		fmt.Println(string(goldory.Value))
	case "purge":
		err = goldory.PurgeAll()
		fmt.Println(string(goldory.Value))
	default:
		log.Fatalf("unsupported client action", axn)
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
}
