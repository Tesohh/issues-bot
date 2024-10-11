package handler

import (
	"issues/db"
	"issues/global"
	"issues/slash"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

var defaultRoles = []*dg.RoleParams{
	{
		Name:        "FEATURE",
		Color:       slash.Ptr(0x00afb9),
		Mentionable: slash.Ptr(true),
	},
	{
		Name:        "FIX",
		Color:       slash.Ptr(0xff8800),
		Mentionable: slash.Ptr(true),
	},
	{
		Name:        "UNITTEST",
		Color:       slash.Ptr(0xf4acb7),
		Mentionable: slash.Ptr(true),
	},
	{
		Name:        "CHORE",
		Color:       slash.Ptr(0xda627d),
		Mentionable: slash.Ptr(true),
	},

	{
		Name:        "CHILL",
		Color:       slash.Ptr(0x0077b6),
		Mentionable: slash.Ptr(true),
	},
	{
		Name:        "IMPORTANT",
		Color:       slash.Ptr(0xffba08),
		Mentionable: slash.Ptr(true),
	},
	{
		Name:        "CRITICAL",
		Color:       slash.Ptr(0xd00000),
		Mentionable: slash.Ptr(true),
	},
}

func AddRoles(s *dg.Session, g *dg.GuildCreate) {
	isNew, err := RegisterGuild(s, g)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if !isNew {
		return
	}

	for i, role := range defaultRoles {
		r, err := s.GuildRoleCreate(g.ID, role)
		if err != nil {
			slog.Error(err.Error())
		}

		var roletype = db.RoleTypeKind
		if i > 3 {
			roletype = db.RoleTypePriority
		}
		global.DB.Create(&db.Role{ID: r.ID, RoleType: roletype, GuildId: g.ID})
	}
}
