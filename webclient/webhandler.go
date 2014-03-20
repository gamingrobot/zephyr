package webclient

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	. "github.com/gamingrobot/zephyr/events"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
)

type WebHandler struct {
	mutex  sync.RWMutex
	client *WebClient

	requestCount uint64
	connections  map[uint64]WebConnection
}

type WebConnection struct {
	webSocket *websocket.Conn
	clientIp  net.Addr
}

func newWebHandler(client *WebClient) *WebHandler {
	return &WebHandler{
		connections: make(map[uint64]WebConnection),
		client:      client,
	}
}

func (w *WebHandler) httpLoop() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory: "public",
	}))
	m.Get("/ws", func(res http.ResponseWriter, req *http.Request) {
		w.webSocketHandler(res, req)
	})
	m.Get("/", func(r render.Render) {
		w.templateIndex(r)
	})
	//go m.Run()
	go w.runServer(m)
	for event := range w.client.webEvents {
		webevent, err := NewWebEvent(event)
		if err != nil {
			fmt.Println("Failed to decode", err)
		} else {
			w.handleWebEvent(webevent)
		}
	}
}

func (w *WebHandler) runServer(m *martini.ClassicMartini) {
	address := "localhost:3000"
	log.Println("Martini started on", address)
	log.Fatal(http.ListenAndServe(address, m))
}

func (w *WebHandler) webSocketHandler(res http.ResponseWriter, req *http.Request) {
	ws, err := websocket.Upgrade(res, req, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(res, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}
	client := ws.RemoteAddr()
	fmt.Println("Got Websocket", client)
	clientId := w.addClient(WebConnection{webSocket: ws, clientIp: client})
	for {
		_, message, err := ws.ReadMessage() //blocking
		if err != nil {
			w.removeClient(clientId)
			return
		} else {
			fmt.Println("Got Event", string(message))
			w.client.webEvents <- string(message)
		}
	}
}

func (w *WebHandler) addClient(connection WebConnection) uint64 {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.requestCount += 1
	w.connections[w.requestCount] = connection
	return w.requestCount
}

func (w *WebHandler) removeClient(id uint64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	delete(w.connections, id)
}

func (w *WebHandler) DispatchEvent(event string) {
	w.mutex.RLock()
	for _, connection := range w.connections {
		err := connection.webSocket.WriteMessage(websocket.TextMessage, []byte(event))
		if err != nil {
			fmt.Println("Websocket write error", err)
		}
	}
	w.mutex.RUnlock()
}

func (w *WebHandler) handleWebEvent(event *WebEvent) {
	fmt.Println(event.Event)
	switch event.Event {
	case "SendMessageEvent":
		w.handleSendMessage(event)
	}
}

func (w *WebHandler) handleSendMessage(event *WebEvent) {
	body := new(SendMessageEvent)
	event.ReadEvent(body)
	steam := w.client.SteamHandler.steam
	steam.Social.SendMessage(body.SteamId, body.ChatEntryType, body.Message)
	fmt.Printf("%+v\n", body)
}
