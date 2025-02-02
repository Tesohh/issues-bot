package command

import (
	"fmt"
	"issues/autolist"
	"issues/db"
	"issues/global"
	"issues/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var List = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "list",
		Description: "list",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "projects",
				Description: "lists projects in current guild",
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "issues",
				Description: "lists issues in current (or specified) project",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionBoolean,
						Name:        "me",
						Description: "only show issues where i'm assigned",
					},
					{
						Type:        dg.ApplicationCommandOptionBoolean,
						Name:        "show_done",
						Description: "shows done and cancelled issues too",
					},
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "priority",
						Description: "only show issues that are of this priority",
					},
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "kind",
						Description: "only show issues that are of this kind",
					},
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "prefix",
						Description: "the project's prefix (eg. ISU, PLB, PYC)",
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
				},
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0]
		options := slash.GetOptionMapRaw(subcommand.Options)

		switch subcommand.Name {
		case "projects":
			var projects []db.Project
			err := global.DB.Find(&projects, "guild_id = ?", i.GuildID).Error
			if err != nil {
				return err
			}

			str := ""
			for _, project := range projects {
				str += fmt.Sprintf("- `%s` %s <#%s>\n", project.Prefix, project.Name, project.IssueChannelID)
			}

			embed := dg.MessageEmbed{
				Title:       fmt.Sprintf("Projects in this guild"),
				Description: str,
			}
			return slash.ReplyWithEmbed(s, i, embed, false)

		case "issues":
			prefixOption, ok := options["prefix"]

			filterMe := false
			{
				meOption, meOk := options["me"]
				filterMe = meOk && meOption.BoolValue()
			}
			showDone := false
			{
				doneOption, doneOk := options["show_done"]
				showDone = doneOk && doneOption.BoolValue()
			}
			var priorityFilter string
			{
				priorityOption, priorityOk := options["priority"]
				if priorityOk {
					priorityFilter = priorityOption.Value.(string)
				}
			}
			var kindFilter string
			{
				kindOption, kindOk := options["kind"]
				if kindOk {
					kindFilter = kindOption.Value.(string)
				}
			}

			var guild db.Guild
			err := global.DB.First(&guild, "id = ?", i.GuildID).Error
			if err != nil {
				return err
			}

			var project db.Project
			if ok {
				prefix := strings.ToUpper(prefixOption.StringValue())
				err := global.DB.
					Preload("Issues").
					First(&project, "prefix = ?", prefix).Error
				if err != nil {
					return err
				}
			} else {
				currentChannel, err := s.Channel(i.ChannelID)
				if err != nil {
					return err
				}
				// get project from channelid or parent channelid
				err = global.DB.
					Preload("Issues").
					First(&project, "issue_channel_id = ? or issue_channel_id = ?", i.ChannelID, currentChannel.ParentID).Error
				if err != nil {
					return err
				}
			}

			filteredIssues := autolist.ApplyFilters(project.Issues, filterMe, i.Member.User.ID, showDone, priorityFilter, kindFilter)

			embedTitle := fmt.Sprintf("Issues for project %s", project.Name)
			embed := autolist.Embed(embedTitle, guild.DefaultPriorityRoleID, filteredIssues)

			return slash.ReplyWithEmbed(s, i, embed, false)
		}

		return nil
	},
}
