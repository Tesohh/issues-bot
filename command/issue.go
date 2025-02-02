package command

import (
	"fmt"
	"issues/autolist"
	"issues/db"
	"issues/global"
	"issues/slash"
	"slices"
	"strings"

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
						Name:        "role",
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
						Name:        "role",
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
		subcommand := i.ApplicationCommandData().Options[0]
		options := slash.GetOptionMapRaw(subcommand.Options)

		lastUpdate := ""

		// get issue from current channel id
		var issue db.Issue
		err := global.DB.Preload("KindRole").Preload("PriorityRole").Find(&issue, "thread_id = ?", i.ChannelID).Error
		if err != nil {
			return err
		}

		if issue.ID == "" {
			return ErrNotInIssue
		}

		var editResultString string

		switch subcommand.Name {
		case "kind", "priority":
			discordRole := options["role"].RoleValue(s, i.GuildID)

			var newRole db.Role
			err = global.DB.Table("roles").Where("id = ?", discordRole.ID).Find(&newRole).Error
			if err != nil {
				return err
			} // no need to check if it's empty on the db, the check role.Roletype != expectedRoleType check already does it

			if subcommand.Name == "kind" {
				if newRole.RoleType != db.RoleTypeKind {
					return ErrRoleIsNotValid
				}
				issue.KindRole = newRole
			} else if subcommand.Name == "priority" {
				if newRole.RoleType != db.RoleTypePriority {
					return ErrRoleIsNotValid
				}
				issue.PriorityRole = newRole
			}

			editResultString = fmt.Sprintf("%s to <@&%s>", subcommand.Name, newRole.ID)
			lastUpdate = fmt.Sprintf("switched %s of %s", subcommand.Name, issue.ID)

		case "assign":
			user := options["person"].UserValue(s)

			assigneeIDs := strings.Split(issue.AssigneeIDs, ",")
			if !slices.Contains(assigneeIDs, user.ID) { // need to add assignee
				issue.AssigneeIDs += "," + user.ID
				editResultString = fmt.Sprintf("assignees and added <@%s>", user.ID)
				lastUpdate = fmt.Sprintf("assigned %s to %s by %s", user.Username, issue.ID, i.Member.User.Username)
			} else { // need to remove assignee
				assigneeIDs = slices.DeleteFunc(assigneeIDs, func(id string) bool {
					return id == user.ID
				})
				issue.AssigneeIDs = strings.Join(assigneeIDs, ",")
				editResultString = fmt.Sprintf("assignees and removed <@%s>", user.ID)
				lastUpdate = fmt.Sprintf("unassigned %s from %s by %s", user.Username, issue.ID, i.Member.User.Username)
			}

		case "rename":
			issue.Title = options["name"].StringValue()
			_, err = s.ChannelEdit(issue.ThreadID, &dg.ChannelEdit{
				Name: issue.ThreadName(),
			})
			if err != nil {
				return err
			}
			editResultString = fmt.Sprintf("the name to %s", issue.Title)
			lastUpdate = fmt.Sprintf("renamed %s by %s", issue.ID, i.Member.User.Username)
		}

		err = global.DB.Save(&issue).Error
		if err != nil {
			return err
		}

		if issue.MessageID == "" {
			editResultString += "\ni was unable to edit the original message for this issue, changes still applied"
		} else {
			_, err = s.ChannelMessageEditEmbed(issue.ThreadID, issue.MessageID, issue.Embed())
			if err != nil {
				return err
			}
		}

		_ = s.InteractionRespond(i, &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> changed %s", i.Member.User.ID, editResultString),
			},
		})

		autolist.Update(nil, issue.ProjectID, nil, i.GuildID, s, lastUpdate)

		return nil
	},
}
