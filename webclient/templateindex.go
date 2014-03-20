package webclient

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	. "github.com/gamingrobot/steamgo/internal"
)

type indexData struct {
	Friends []Friend
}

type Friend struct {
	SteamId    uint64
	Name       string
	AvatarRoot string
	Avatar     string
	State      EPersonaState
}

func (w *WebHandler) templateIndex(r render.Render) {
	index := indexData{}
	steam := w.client.SteamHandler.steam
	for _, f := range steam.Social.Friends.GetCopy() {
		avatar := f.Avatar
		if !ValidAvatar(avatar) {
			avatar = DefaultAvatar
		}
		fmt.Println(f.SteamId.ToUint64(), f.Name)
		friend := Friend{
			SteamId:    f.SteamId.ToUint64(),
			Name:       f.Name,
			AvatarRoot: avatar[0:2],
			Avatar:     avatar,
			State:      f.PersonaState,
		}
		index.Friends = append(index.Friends, friend)
	}
	r.HTML(200, "index", index)
}
