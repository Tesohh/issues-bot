package handler

import (
	"issues/db"
	"issues/global"
	"issues/slash"

	dg "github.com/bwmarrin/discordgo"
)

var kindRoles = []*dg.RoleParams{
	{
		Name:        "GENERIC",
		Color:       slash.Ptr(0xfffffc),
		Mentionable: slash.Ptr(true),
	},
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
}
var priorityRoles = []*dg.RoleParams{
	{
		Name:        "LOW",
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

func addRole(s *dg.Session, g *dg.GuildCreate, role *dg.RoleParams, roleType db.RoleType) (db.Role, error) {
	r, err := s.GuildRoleCreate(g.ID, role)
	if err != nil {
		return db.Role{}, err
	}

	dbrole := db.Role{ID: r.ID, RoleType: roleType, GuildId: g.ID}
	result := global.DB.Create(&dbrole)
	return dbrole, result.Error
}

func AddRoles(s *dg.Session, g *dg.GuildCreate) error {
	var defaultKindRole db.Role
	for i, role := range kindRoles {
		r, err := addRole(s, g, role, db.RoleTypeKind)
		if err != nil {
			return err
		}
		if i == 0 {
			defaultKindRole = r
		}
	}

	result := global.DB.
		Table("guilds").
		Where("id = ?", g.ID).
		Updates(db.Guild{DefaultKindRoleID: defaultKindRole.ID})

	if result.Error != nil {
		return result.Error
	}

	var defaultPriorityRole db.Role
	for i, role := range priorityRoles {
		r, err := addRole(s, g, role, db.RoleTypePriority)
		if err != nil {
			return err
		}
		if i == 0 {
			defaultPriorityRole = r
		}
	}
	result = global.DB.
		Table("guilds").
		Where("id = ?", g.ID).
		Updates(db.Guild{DefaultPriorityRoleID: defaultPriorityRole.ID})

	return result.Error
}
