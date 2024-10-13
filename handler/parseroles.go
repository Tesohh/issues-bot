package handler

import (
	"issues/db"
	"issues/global"
	"slices"

	dg "github.com/bwmarrin/discordgo"
)

func ParseRoles(s *dg.Session, guildID string, roleIDs []string) (*db.Role, *db.Role, error) {
	var guild db.Guild
	result := global.DB.
		Table("guilds").
		Preload("DefaultKindRole").
		Preload("DefaultPriorityRole").
		Find(&guild, "id = ?", guildID)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var roles []db.Role
	result = global.DB.Table("roles").Find(&roles, "guild_id = ?", guildID)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	kindRoles := []db.Role{}
	priorityRoles := []db.Role{}

	for _, role := range roles {
		if slices.Contains(roleIDs, role.ID) {
			if role.RoleType == db.RoleTypeKind {
				kindRoles = append(kindRoles, role)
			} else if role.RoleType == db.RoleTypePriority {
				priorityRoles = append(priorityRoles, role)
			}
		}
	}

	kindRole := guild.DefaultKindRole
	priorityRole := guild.DefaultPriorityRole

	if len(kindRoles) > 0 {
		kindRole = kindRoles[0]
	}
	if len(priorityRoles) > 0 {
		priorityRole = priorityRoles[0]
	}

	return &kindRole, &priorityRole, nil
}
