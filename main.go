package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/gamingrobot/steamgo"
	. "github.com/gamingrobot/steamgo/internal"
	. "github.com/gamingrobot/steamgo/steamid"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
)

type Test struct {
	Testing string
}

type WebConn struct {
	webSocket *websocket.Conn
	clientIp  net.Addr
}

type Router struct {
	mutex        sync.RWMutex
	requestCount uint64
	connections  map[uint64]WebConn
}

type Event struct {
	Name  string
	Event interface{}
}

var router Router
var logger *log.Logger
var lastnum int

func main() {
	logger = log.New(os.Stdout, "[debug] ", log.Lshortfile|log.Lmicroseconds)
	runtime.GOMAXPROCS(4)
	webevents := make(chan string, 20)
	steamevents := make(chan string, 100)
	go startRouter(steamevents)
	go startHttp(webevents)
	startSteam(webevents, steamevents)
}

func startRouter(steamevents <-chan string) {
	logger.Println("Router Started")
	router.connections = make(map[uint64]WebConn)
	for event := range steamevents {
		router.mutex.RLock()
		//logger.Println("Number of connections", len(router.connections))
		for _, connection := range router.connections {
			err := connection.webSocket.WriteMessage(websocket.TextMessage, []byte(event))
			if err != nil {
				logger.Println("Websocket write error:", err)
			}
		}
		router.mutex.RUnlock()
	}
}

func WebSocketHandler(res http.ResponseWriter, req *http.Request, webevents chan<- string) {
	ws, err := websocket.Upgrade(res, req, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(res, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		logger.Println(err)
		return
	}
	client := ws.RemoteAddr()
	sockCli := WebConn{webSocket: ws, clientIp: client}
	clientId := addClient(sockCli)

	for {
		_, message, err := ws.ReadMessage() //blocking
		if err != nil {
			removeClient(clientId)
			return
		} else {
			webevents <- string(message)
		}
	}
}

func startHttp(webevents chan<- string) {
	m := martini.Classic()
	logger.Println("Martini Started")
	m.Get("/ws", func(res http.ResponseWriter, req *http.Request) {
		WebSocketHandler(res, req, webevents)
	})
	m.Run()
}

func startSteam(webevents <-chan string, steamevents chan<- string) {
	logger.Println("Steam Started")
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	login := steamgo.LogOnDetails{}
	decoder.Decode(&login)
	client := steamgo.NewClient()
	server := client.ConnectNorthAmerica()
	logger.Println("Connected to server:", server)
	for {
		select {
		case jsonevent := <-webevents:
			logger.Println("WebEvent", jsonevent)
			webevent := Event{}
			err := json.Unmarshal([]byte(jsonevent), &webevent)
			if err != nil {
				logger.Println("Failed to decode event", err)
			} else {
				e := webevent.Event
				event := e.(map[string]interface{})
				switch webevent.Name {
				case "SendChatMessage":
					logger.Printf("%+v\n", event)
					steamid, _ := strconv.ParseUint(event["SteamId"].(string), 10, 64)
					chatentry := int32(event["ChatEntryType"].(float64))
					client.Social.SendChatMessage(
						SteamId(steamid),
						EChatEntryType(chatentry),
						event["Message"].(string))
				}
			}
		case steamevent := <-client.Events():
			switch e := steamevent.(type) {
			case steamgo.ConnectedEvent:
				client.Auth.LogOn(login)
			case steamgo.LoggedOnEvent:
				client.Social.SetPersonaState(EPersonaState_Online)
			case steamgo.FatalError:
				client.Connect() // please do some real error handling here
				logger.Print("FatalError")
				logger.Print(e)
			case error:
				logger.Print(e)
			}
			outevent, err := json.Marshal(Event{
				Name:  reflect.TypeOf(steamevent).Name(),
				Event: steamevent,
			})
			if err != nil {
				logger.Println("Failed to encode event", err)
			} else {
				steamevents <- string(outevent)
			}
		}
	}
}

func addClient(connection WebConn) uint64 {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.requestCount += 1
	router.connections[router.requestCount] = connection
	return router.requestCount
}

func removeClient(id uint64) {
	logger.Println("Removing Client ", id)
	router.mutex.Lock()
	defer router.mutex.Unlock()
	delete(router.connections, id)
}
