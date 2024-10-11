package handler

import (
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

func GuildJoinHandler(s *dg.Session, g *dg.GuildCreate) {
	isNew, err := RegisterGuild(s, g)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if !isNew {
		return
	}

	err = AddRoles(s, g)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
