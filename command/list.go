package command

import (
	"fmt"
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

			issueStrings := make([]string, 4)

			for _, issue := range project.Issues {
				issueStrings[issue.IssueStatus] += fmt.Sprintf("- <#%s>\n", issue.ThreadID)
			}

			for i := range issueStrings {
				if issueStrings[i] == "" {
					issueStrings[i] = "nothing here..."
				} else {
					issueStrings[i] = strings.TrimRight(issueStrings[i], "\n")
				}
			}

			description := fmt.Sprintf("**Todo**\n%s\n**Doing**\n%s\n**Done**\n%s\n**Cancelled**\n%s", issueStrings[0], issueStrings[1], issueStrings[2], issueStrings[3])

			embed := dg.MessageEmbed{
				Title:       fmt.Sprintf("Issues in project %s", project.Name),
				Description: description,
			}
			return slash.ReplyWithEmbed(s, i, embed, false)
		}

		return nil
	},
}
