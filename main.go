package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/gamingrobot/steamgo"
	. "github.com/gamingrobot/steamgo/internal"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type Router struct {
	mutex        sync.RWMutex
	requestcount uint64
	channels     map[uint64]chan string
}

var router Router
var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "[debug] ", log.Lshortfile|log.Lmicroseconds)
	runtime.GOMAXPROCS(4)
	webevents := make(chan string)
	steamevents := make(chan string, 100)
	go startRouter(steamevents)
	go startHttp(webevents)
	startSteam(webevents, steamevents)
}

func startRouter(steamevents <-chan string) {
	logger.Println("startRouter")

	router.channels = make(map[uint64]chan string)
	for event := range steamevents {
		logger.Println("New Event")

		router.mutex.Lock()
		logger.Println("Number of channels", len(router.channels))
		for k, channel := range router.channels {
			logger.Println("Sending event to ", k)

			logger.Println("Event Written", event)
			channel <- event
		}
		router.mutex.Unlock()
		logger.Println("Event routed")

	}
}

func startHttp(webevents chan<- string) {
	m := martini.Classic()
	logger.Println("Martini")

	m.Get("/poll", func(params martini.Params) (int, string) {
		logger.Println("Poll HTTP")
		timeout, err := strconv.Atoi(params["timeout"])
		if err != nil {
			timeout = 60000 //default 60 seconds
		}
		timeout -= 100 //remove 100 ms from the timout
		tempchan := make(chan string, 100)
		requestnum := addChannel(tempchan)

		logger.Println("XHR ID is", requestnum)
		defer removeChannel(requestnum)
		//defer close(tempchan)
		select {
		case event := <-tempchan:
			logger.Println("Event Sent", event)

			logger.Println("I am ", requestnum, " and I got a event")
			return 200, event
		case <-time.After(time.Duration(30000) * time.Millisecond):
			return 400, ""
		}
	})
	m.Get("/send/:message", func(params martini.Params) string {
		webevents <- params["message"]
		return `{"err": false}`
	})
	m.Run()
}

func startSteam(webevents <-chan string, steamevents chan<- string) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	login := steamgo.LogOnDetails{}
	decoder.Decode(&login)
	client := steamgo.NewClient()
	server := client.ConnectNorthAmerica()
	logger.Println("Connected to server:", server)
	for {
		select {
		case webevent := <-webevents:
			logger.Println("WebEvent", webevent)
		case steamevent := <-client.Events():
			switch e := steamevent.(type) {
			case *steamgo.ConnectedEvent:
				client.Auth.LogOn(&login)
			case *steamgo.LoggedOnEvent:
				client.Social.SetPersonaState(EPersonaState_Online)
			case steamgo.FatalError:
				client.Connect() // please do some real error handling here
				logger.Print("FatalError")
				logger.Print(e)
			case error:
				logger.Print(e)
			}
			m, err := json.Marshal(steamevent)
			if err == nil {
				steamevents <- string(m)
			}
		}
	}
}

func addChannel(newchan chan string) uint64 {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.requestcount += 1
	router.channels[router.requestcount] = newchan
	return router.requestcount
}

func removeChannel(id uint64) {
	logger.Println("Removing Chan ", id)

	router.mutex.Lock()
	defer router.mutex.Unlock()
	delete(router.channels, id)
}
