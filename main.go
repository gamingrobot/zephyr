package main

import (
	"encoding/json"
	"github.com/gamingrobot/steamgo"
	"github.com/gamingrobot/zephyr/webclient"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	//load login details
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	login := steamgo.LogOnDetails{}
	decoder.Decode(&login)
	webclient := webclient.NewWebClient()
	webclient.Start(login)
}
