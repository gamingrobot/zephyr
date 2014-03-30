package webclient

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	. "github.com/gamingrobot/steamgo/internal"
	"github.com/gamingrobot/steamgo/socialcache"
	"strconv"
)

type indexData struct {
	User    User
	Friends []Friend
	Groups  []Group
	Chats   []Chat
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

type Group struct {
	SteamId    uint64
	Name       string
	AvatarRoot string
	Avatar     string
	StateText  string
}

type Chat struct {
	SteamId     uint64
	GroupId     uint64
	ChatMembers []ChatMember
}

type ChatMember struct {
	SteamId uint64
	Rank    string
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
	for _, g := range steam.Social.Groups.GetCopy() {
		avatar := g.Avatar
		if !ValidAvatar(avatar) {
			avatar = DefaultAvatar
		}
		fmt.Println(g.SteamId.ToUint64(), g.Name)
		group := Group{
			SteamId:    g.SteamId.ToUint64(),
			Name:       g.Name,
			AvatarRoot: avatar[0:2],
			Avatar:     avatar,
			StateText:  strconv.FormatUint(uint64(g.MemberChattingCount), 10),
		}
		index.Groups = append(index.Groups, group)
	}
	for _, c := range steam.Social.Chats.GetCopy() {
		chat := Chat{
			SteamId: c.SteamId.ToUint64(),
			GroupId: c.GroupId.ToUint64(),
		}
		for _, cm := range c.ChatMembers {
			chatmember := ChatMember{
				SteamId: cm.SteamId.ToUint64(),
				Rank:    getRank(cm),
			}
			chat.ChatMembers = append(chat.ChatMembers, chatmember)
		}
		index.Chats = append(index.Chats, chat)
	}
	r.HTML(200, "index", index)
}

func getRank(c socialcache.ChatMember) string {
	perm := c.ClanPermissions
	switch perm {
	case EClanPermission_Nobody:
		return "none"
	case EClanPermission_Owner:
		return "admin"
	case EClanPermission_Officer:
		return "admin"
	case EClanPermission_OwnerAndOfficer:
		return "admin"
	case EClanPermission_Member:
		return "member"
	case EClanPermission_Moderator:
		return "mod"
	case EClanPermission_OwnerOfficerModerator:
		return "admin"
	default:
		return "none"
	}
}

func getState(f socialcache.Friend) (string, string) {
	if f.GameId != 0 {
		return "ingame", f.GameName
	}
	state := f.PersonaState
	switch state {
	case EPersonaState_Away:
		return "away", "Away"
	case EPersonaState_Busy:
		return "busy", "Busy"
	case EPersonaState_Offline:
		return "offline", "Offline"
	case EPersonaState_LookingToPlay:
		return "lookingtoplay", "Looking to Play"
	case EPersonaState_LookingToTrade:
		return "lookingtotrade", "Looking to Trade"
	case EPersonaState_Online:
		return "online", "Online"
	case EPersonaState_Snooze:
		return "snooze", "Snooze"
	default:
		return "offline", "Offline"
	}
}
