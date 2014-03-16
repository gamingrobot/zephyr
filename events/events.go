package events

import (
	"encoding/json"
	. "github.com/gamingrobot/steamgo/internal"
	. "github.com/gamingrobot/steamgo/steamid"
)

type Event string

type SendMessageEvent struct {
	SteamId       SteamId
	ChatEntryType EChatEntryType
	Message       string
}

type WebEvent struct {
	Event Event
	Body  json.RawMessage
}

type SteamEvent struct {
	Event Event
	Body  interface{}
}
