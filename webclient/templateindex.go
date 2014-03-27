package webclient

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	. "github.com/gamingrobot/steamgo/internal"
	"github.com/gamingrobot/steamgo/socialcache"
)

type indexData struct {
	User    User
	Friends []Friend
}

type User struct {
	SteamId    uint64
	Name       string
	AvatarRoot string
	Avatar     string
}

type Friend struct {
	SteamId    uint64
	Name       string
	AvatarRoot string
	Avatar     string
	State      string
	StateText  string
}

func (w *WebHandler) templateIndex(r render.Render) {
	index := indexData{}
	steam := w.client.steamClient
	uavatar := steam.Social.GetAvatar()
	if !ValidAvatar(uavatar) {
		uavatar = DefaultAvatar
	}
	index.User = User{
		SteamId:    uint64(steam.SteamId()),
		Name:       steam.Social.GetPersonaName(),
		AvatarRoot: uavatar[0:2],
		Avatar:     uavatar,
	}
	for _, f := range steam.Social.Friends.GetCopy() {
		avatar := f.Avatar
		if !ValidAvatar(avatar) {
			avatar = DefaultAvatar
		}
		fmt.Println(f.SteamId.ToUint64(), f.Name)
		state, stateText := getState(f)
		friend := Friend{
			SteamId:    f.SteamId.ToUint64(),
			Name:       f.Name,
			AvatarRoot: avatar[0:2],
			Avatar:     avatar,
			State:      state,
			StateText:  stateText,
		}
		index.Friends = append(index.Friends, friend)
	}
	r.HTML(200, "index", index)
}

func getState(f socialcache.Friend) (string, string) {
	if f.GameId != 0 {
		return "ingame", f.GameName
	}
	state := f.PersonaState
	if state == EPersonaState_Away {
		return "away", "Away"
	} else if state == EPersonaState_Busy {
		return "busy", "Busy"
	} else if state == EPersonaState_Offline {
		return "offline", "Offline"
	} else if state == EPersonaState_LookingToPlay {
		return "lookingtoplay", "Looking to Play"
	} else if state == EPersonaState_LookingToTrade {
		return "lookingtotrade", "Looking to Trade"
	} else if state == EPersonaState_Online {
		return "online", "Online"
	} else if state == EPersonaState_Snooze {
		return "snooze", "Snooze"
	}
	return "offline", "Offline"

}
