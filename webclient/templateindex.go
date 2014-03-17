package webclient

import (
	"github.com/codegangsta/martini-contrib/render"
)

type indexData struct {
	Friends []Friend
}

type Friend struct {
	SteamId    uint64
	Name       string
	AvatarRoot string
	Avatar     string
	State      string
}

func (w *WebHandler) templateIndex(r render.Render) {
	index := indexData{}
	steam := w.client.SteamHandler.steam
	for _, f := range steam.Social.Friends.GetCopy() {
		friend := Friend{
			SteamId:    f.SteamId.ToUint64(),
			Name:       f.Name,
			AvatarRoot: f.Avatar[0:2],
			Avatar:     f.Avatar,
			State:      string(f.PersonaState),
		}
		index.Friends = append(index.Friends, friend)
	}
	r.HTML(200, "index", index)
}
