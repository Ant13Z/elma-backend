package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"./daemon"
	"./db"
)

type Config struct {
	Db db.Config
	Server daemon.Config
}

func getConfig() Config{
	var arr Config
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(file, &arr)
	return arr
}

func startServer(arr Config) {
	daemon.Run(arr.Server, arr.Db)
}

func main() {
	startServer(getConfig())
}