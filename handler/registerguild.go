package handler

import (
	"issues/db"
	"issues/global"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

// note this is not a handler, this is for AddRoles
func RegisterGuild(s *dg.Session, g *dg.GuildCreate) (bool, error) {
	slog.Info("checking guild..")
	guild := db.Guild{ID: g.ID}
	result := global.DB.FirstOrCreate(&guild, guild)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected > 0 {
		slog.Info("registered guild", "id", g.ID)
		return true, nil
	} else {
		slog.Info("connected to guild", "id", g.ID)
		return false, nil
	}
}
