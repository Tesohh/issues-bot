package command

import (
	"issues/db"
	"issues/global"
	"issues/slash"

	dg "github.com/bwmarrin/discordgo"
)

var RegisterRole = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "registerrole",
		Description: "register a new role",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "kind",
				Description: "registers a new kind role",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "the kind role",
						Required:    true,
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "priority",
				Description: "registers a new priority role",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "the priority role",
						Required:    true,
					},
				},
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0]
		options := slash.GetOptionMapRaw(subcommand.Options)

		var roleType db.RoleType
		if subcommand.Name == "kind" {
			roleType = db.RoleTypeKind
		} else if subcommand.Name == "priority" {
			roleType = db.RoleTypePriority
		}

		role := options["role"].RoleValue(s, i.GuildID)

		var exists bool
		err := global.DB.Table("roles").
			Select("count(*) > 0").
			Where("id = ?", role.ID).
			Find(&exists).
			Error
		if err != nil {
			return err
		} else if exists {
			return ErrRoleAlreadyExists
		}

		dbRole := db.Role{
			ID:       role.ID,
			RoleType: roleType,
			GuildId:  i.GuildID,
		}
		err = global.DB.Save(&dbRole).Error
		if err != nil {
			return err
		}

		embed := dg.MessageEmbed{
			Title: "Registered role",
		}
		return slash.ReplyWithEmbed(s, i, embed, true)
	},
}
