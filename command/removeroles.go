package command

import (
	"issues/db"
	"issues/global"
	"issues/slash"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var RemoveRoles = slash.Command{
	ApplicationCommand: discordgo.ApplicationCommand{
		Name:        "removeroles",
		Description: "Removes all roles registered by me",
	},
	Func: func(s *discordgo.Session, i *discordgo.Interaction) error {
		slog.Info("REMOVING ROLES...")
		roles := []db.Role{}
		global.DB.Find(&roles)

		for _, role := range roles {
			s.GuildRoleDelete(i.GuildID, role.ID)
			global.DB.Delete(&role)
		}

		return nil
	},
}
