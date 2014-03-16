package webclient

import (
	"fmt"
	"github.com/gamingrobot/steamgo"
	. "github.com/gamingrobot/steamgo/internal"
)

type SteamHandler struct {
	client *WebClient
}

func newSteamHandler(client *WebClient) *SteamHandler {
	return &SteamHandler{
		client: client,
	}
}

func (s *SteamHandler) steamLoop(login steamgo.LogOnDetails) {
	client := steamgo.NewClient()
	server := client.ConnectNorthAmerica()
	fmt.Println("Connected to server:", server)
	for event := range client.Events() {
		switch e := event.(type) {
		case steamgo.ConnectedEvent:
			client.Auth.LogOn(login)
		case steamgo.FatalError:
			client.Connect() // please do some real error handling here
			fmt.Print("FatalError", e)
		case error:
			fmt.Println(e)
		default:
			s.handleSteamEvent(event, client)
		}
	}
}

func (s *SteamHandler) handleSteamEvent(event interface{}, client *steamgo.Client) {
	switch event.(type) {
	case steamgo.LoggedOnEvent:
		client.Social.SetPersonaState(EPersonaState_Online)
	}
}
