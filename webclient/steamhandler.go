package webclient

import (
	"github.com/gamingrobot/steamgo"
	. "github.com/gamingrobot/steamgo/internal"
	. "github.com/gamingrobot/zephyr/events"
	"log"
)

type SteamHandler struct {
	client *WebClient
	steam  *steamgo.Client
}

func newSteamHandler(client *WebClient) *SteamHandler {
	steam := steamgo.NewClient()
	server := steam.ConnectNorthAmerica()
	log.Println("Connecting to steam server:", server)
	return &SteamHandler{
		client: client,
		steam:  steam,
	}
}

func (s *SteamHandler) steamLoop(login steamgo.LogOnDetails) {
	for event := range s.steam.Events() {
		switch e := event.(type) { //Events that should *not* be passed to web
		case steamgo.ConnectedEvent:
			log.Println("Connected to steam")
			s.steam.Auth.LogOn(login)
		case steamgo.LoggedOnEvent:
			log.Println("Logged on steam as", login.Username)
		case steamgo.LoggedOffEvent:
			log.Println("Logged off steam")
		case steamgo.DisconnectedEvent:
			log.Println("Disconnected to steam")
		case steamgo.MachineAuthUpdateEvent:
		case steamgo.LoginKeyEvent:
		case steamgo.FatalError:
			s.steam.Connect() // please do some real error handling here
			log.Print("FatalError", e)
		case error:
			log.Println(e)
		default:
			s.handleSteamEvent(event)
		}
	}
}

func (s *SteamHandler) handleSteamEvent(event interface{}) {
	switch event.(type) { //Events that should be passed to web
	case steamgo.LoggedOnEvent:
		s.steam.Social.SetPersonaState(EPersonaState_Online)
	}
	steamevent, err := EncodeEvent(event)
	if err != nil {
		log.Println("Failed to encode", err)
	} else {
		s.client.steamEvents <- steamevent
	}
}
