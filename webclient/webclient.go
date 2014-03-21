package webclient

import (
	"github.com/gamingrobot/steamgo"
)

type WebClient struct {
	WebHandler   *WebHandler
	SteamHandler *SteamHandler
	steamClient  *steamgo.Client

	webEvents   chan string
	steamEvents chan string
}

func NewWebClient() *WebClient {
	client := &WebClient{
		webEvents:   make(chan string, 20),
		steamEvents: make(chan string, 100),
	}
	client.WebHandler = newWebHandler(client)
	client.SteamHandler = newSteamHandler(client)
	client.steamClient = client.SteamHandler.steam
	return client
}

func (c *WebClient) Start(login steamgo.LogOnDetails) {
	go c.WebHandler.httpLoop()
	go c.SteamHandler.steamLoop(login)
	c.routerLoop()

}

func (c *WebClient) routerLoop() {
	for event := range c.steamEvents {
		c.WebHandler.DispatchEvent(event)
	}
}
