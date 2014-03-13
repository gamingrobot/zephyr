package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/gamingrobot/steamgo"
	. "github.com/gamingrobot/steamgo/internal"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
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
var lastnum int

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
		logger.Println("New Event in Router")
		router.mutex.Lock()
		logger.Println("Number of channels", len(router.channels))
		if len(router.channels) < lastnum {
			time.Sleep(time.Duration(500) * time.Millisecond)
			logger.Println("WOAH WHAT HAPPEND")
		}
		for k, channel := range router.channels {
			logger.Println("Sending event to ", k)

			logger.Println("Event Written to Martini", event)
			channel <- event
		}
		lastnum = len(router.channels)
		router.mutex.Unlock()
		logger.Println("Event routed")

	}
}

func PollHandler(rw http.ResponseWriter, req *http.Request, params martini.Params) (int, string) {
	logger.Println("Poll HTTP")
	timeout, err := strconv.ParseInt(req.URL.Query().Get("timeout"), 10, 64)
	if err != nil {
		timeout = 30000 //default 30 seconds
	}
	timeout -= 1000 //remove 1000 ms from the timout
	tempchan := make(chan string, 1)
	requestnum := addChannel(tempchan)

	logger.Println("XHR ID is", requestnum)
	defer removeChannel(requestnum)
	//defer close(tempchan)
	select {
	case event := <-tempchan:
		logger.Println("Event Sent to Client", event)

		logger.Println("I am ", requestnum, " and I got a event")
		return 200, event
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		return 400, ""
	}
}

func startHttp(webevents chan<- string) {
	m := martini.Classic()
	logger.Println("Martini")

	m.Get("/poll", PollHandler)
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
			logger.Println("Event Got from Steam", steamevent)
			m, err := json.Marshal(steamevent)
			if err != nil {
				logger.Println("FAILED TO ENCODE THIS THING")
			} else {
				steamevents <- string(m)
				logger.Println("Event Sent to Router", string(m))
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
