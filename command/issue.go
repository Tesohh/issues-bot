package command

import (
	"fmt"
	"issues/db"
	"issues/global"
	"issues/slash"

	dg "github.com/bwmarrin/discordgo"
)

var Issue = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "issue",
		Description: "manage issues",
		Options: []*dg.ApplicationCommandOption{
			// by design you cant:
			// - delete (you would set as cancelled)
			// - change the description (you would just talk about it in the channel)
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "kind",
				Description: "edits the kind of this issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "kind",
						Description: "the kind role",
						Required:    true,
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "priority",
				Description: "edits the priority of this issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "priority",
						Description: "the priority role",
						Required:    true,
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "assign",
				Description: "assign this person, remove if they were already assigned",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionUser,
						Name:        "person",
						Description: "person to assign/unassign",
						Required:    true,
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "rename",
				Description: "rename the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "name",
						Description: "new name",
						Required:    true,
					},
				},
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0].Name
		options := slash.GetOptionMap(i)

		// get issue from current channel id
		var issue db.Issue
		err := global.DB.Find(&issue, "thread_id = ?", i.ChannelID).Error
		if err != nil {
			return err
		}

		if issue.ID == "" {
			return ErrNotInIssue
		}

		var editResultString string

		switch subcommand {
		case "kind", "priority":
			discordRole := options[subcommand].RoleValue(s, i.GuildID)

			var newRole db.Role
			err = global.DB.Table("roles").Where("id = ?", discordRole.ID).Find(&newRole).Error
			if err != nil {
				return err
			} // no need to check if it's empty on the db, the check role.Roletype != expectedRoleType check already does it

			var expectedRoleType db.RoleType
			if subcommand == "kind" {
				expectedRoleType = db.RoleTypeKind
			} else if subcommand == "priority" {
				expectedRoleType = db.RoleTypePriority
			}

			if newRole.RoleType != expectedRoleType {
				return ErrRoleIsNotValid
			}

			for i, role := range issue.Roles {
				if role.RoleType == expectedRoleType {
					issue.Roles[i] = newRole
				}
			}

			err = global.DB.Save(issue).Error
			if err != nil {
				return err
			}
			editResultString = fmt.Sprintf("%s to <@&%s>", subcommand, newRole.ID)

		case "assign":
		case "rename":
		}

		if issue.MessageID == "" {
			editResultString += "\ni was unable to edit the original message for this issue, changes still applied"
		}

		_, _ = s.ChannelMessageSend(issue.ThreadID, fmt.Sprintf("<@%s> changed %s", i.User.ID, editResultString))

		_, err = s.ChannelMessageEditEmbed(issue.ThreadID, issue.MessageID, issue.Embed())
		if err != nil {
			return err
		}

		return nil
	},
}
