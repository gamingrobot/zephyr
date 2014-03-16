package events

import (
	"encoding/json"
	"reflect"
)

func (e *WebEvent) ReadEvent(body interface{}) {
	json.Unmarshal([]byte(e.Body), &body)
}

func NewWebEvent(in string) (*WebEvent, error) {
	event := new(WebEvent)
	err := json.Unmarshal([]byte(in), event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func EncodeEvent(event interface{}) (string, error) {
	encoded, err := json.Marshal(SteamEvent{
		Event: Event(reflect.TypeOf(event).Name()),
		Body:  event,
	})
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
