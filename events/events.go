package events

import (
	"encoding/json"
	. "github.com/gamingrobot/steamgo/internal"
	. "github.com/gamingrobot/steamgo/steamid"
	"strconv"
)

type Event string

type WebEvent struct {
	Event Event
	Body  json.RawMessage
}

type SteamEvent struct {
	Event Event
	Body  interface{}
}

type SendMessageEvent struct {
	SteamId       JSteamId
	ChatEntryType EChatEntryType
	Message       string
}

//Helper types because Javascript is dumb
type JSteamId string
type JUint64 string
type JInt64 string

func (s JSteamId) Convert() SteamId {
	id, _ := strconv.ParseUint(string(s), 10, 64)
	return SteamId(id)
}

func (s JUint64) Convert() uint64 {
	num, _ := strconv.ParseUint(string(s), 10, 64)
	return num
}

func (s JInt64) Convert() int64 {
	num, _ := strconv.ParseInt(string(s), 10, 64)
	return num
}
