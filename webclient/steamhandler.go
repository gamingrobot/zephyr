package webclient

import (
	"fmt"
	"github.com/gamingrobot/steamgo"
	. "github.com/gamingrobot/steamgo/internal"
)

type SteamHandler struct {
	client *WebClient
	steam  *steamgo.Client
}

func newSteamHandler(client *WebClient) *SteamHandler {
	steam := steamgo.NewClient()
	server := steam.ConnectNorthAmerica()
	fmt.Println("Connected to server:", server)
	return &SteamHandler{
		client: client,
		steam:  steam,
	}
}

func (s *SteamHandler) steamLoop(login steamgo.LogOnDetails) {
	for event := range s.steam.Events() {
		switch e := event.(type) {
		case steamgo.ConnectedEvent:
			s.steam.Auth.LogOn(login)
		case steamgo.FatalError:
			s.steam.Connect() // please do some real error handling here
			fmt.Print("FatalError", e)
		case error:
			fmt.Println(e)
		default:
			s.handleSteamEvent(event)
		}
	}
}

func (s *SteamHandler) handleSteamEvent(event interface{}) {
	switch event.(type) {
	case steamgo.LoggedOnEvent:
		s.steam.Social.SetPersonaState(EPersonaState_Online)
	}
}
