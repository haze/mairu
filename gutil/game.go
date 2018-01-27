package game

import (
	"github.com/bwmarrin/discordgo"
)

// Valid GameType Values
const (
	TypeGame discordgo.GameType = iota
	TypeStreaming
	TypeListening
	TypeWatching
)

// UpdateStatusSpecial .. does cool shit
func UpdateStatusSpecial(sesh *discordgo.Session, afk bool, game string, t discordgo.GameType) {
	idle := 0
	status := discordgo.UpdateStatusData{
		IdleSince: &idle,
		Game: &discordgo.Game{
			Name: game,
			Type: t,
		},
		AFK:    afk,
		Status: "online",
	}
	sesh.UpdateStatusComplex(status)
}
