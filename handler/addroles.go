package handler

import (
	"issues/slash"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

var defaultRoles = []*dg.RoleParams{
	{
		Name:        "ISSUESTEST",
		Color:       slash.Ptr(0x2a9d8f),
		Mentionable: slash.Ptr(true),
	},
	// {
	// 	Name:        "FIX",
	// 	Color:       slash.Ptr(0xc1121f),
	// 	Mentionable: slash.Ptr(true),
	// },
	// {
	// 	Name:        "UNITTEST",
	// 	Color:       slash.Ptr(0x3d348b),
	// 	Mentionable: slash.Ptr(true),
	// },
	//
	// {
	// 	Name:        "CHILL",
	// 	Color:       slash.Ptr(0x0077b6),
	// 	Mentionable: slash.Ptr(true),
	// },
	// {
	// 	Name:        "IMPORTANT",
	// 	Color:       slash.Ptr(0xffba08),
	// 	Mentionable: slash.Ptr(true),
	// },
	// {
	// 	Name:        "CRITICAL",
	// 	Color:       slash.Ptr(0xd00000),
	// 	Mentionable: slash.Ptr(true),
	// },
}

func AddRoles(s *dg.Session, g *dg.GuildCreate) {
	slog.Info("joined guild.")
	slog.Warn("no roles will be created as we still can't check if we are in a new guild") // TEMP:
	// for _, role := range defaultRoles {
	// 	_, err := s.GuildRoleCreate(g.ID, role)
	// 	if err != nil {
	// 		slog.Error(err.Error())
	// 	}
	// }
}
