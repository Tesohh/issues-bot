package command

import (
	"errors"
	"fmt"
	"issues/autolist"
	"issues/db"
	"issues/global"
	"issues/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var Project = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "project",
		Description: "manage projects",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "new",
				Description: "create a new project",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "prefix",
						Description: "the project's prefix (eg. ISU, PLB, PYC)",
						Required:    true,
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "name",
						Description: "the project's name",
						Required:    true,
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "rename",
				Description: "rename a project",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "prefix",
						Description: "the project's prefix (eg. ISU, PLB, PYC)",
						Required:    true,
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "name",
						Description: "the project's name",
						Required:    true,
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "delete a project",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "prefix",
						Description: "the project's prefix (eg. ISU, PLB, PYC)",
						Required:    true,
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
					{
						Type:        dg.ApplicationCommandOptionBoolean,
						Name:        "confirm",
						Description: "are you really sure?",
						Required:    true,
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "resetlistmsg",
				Description: "in case shit hits the fan, you can reset the list message by sending a new one",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "prefix",
						Description: "the project's prefix (eg. ISU, PLB, PYC)",
						Required:    true,
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
		prefix := strings.ToUpper(options["prefix"].StringValue())

		var project db.Project
		err := global.DB.Find(&project, "prefix = ?", prefix).Error
		if err != nil {
			return err
		}

		switch subcommand.Name {
		case "new":
			name := options["name"].StringValue()
			category, err := s.GuildChannelCreate(i.GuildID, name, dg.ChannelTypeGuildCategory)
			if err != nil {
				return err
			}

			issuesChannelName := fmt.Sprintf("%s-issues", prefix)
			issuesChannel, err := s.GuildChannelCreateComplex(i.GuildID, dg.GuildChannelCreateData{
				Name:     issuesChannelName,
				Type:     dg.ChannelTypeGuildText,
				ParentID: category.ID,
			})
			if err != nil {
				return err
			}

			// TODO: Send the AutoList message and set id
			embedTitle := fmt.Sprintf("AutoList™️ for %s", name)
			alMsg, err := s.ChannelMessageSendEmbed(issuesChannel.ID, slash.Ptr(autolist.Embed(embedTitle, "", []db.Issue{})))
			if err != nil {
				return err
			}

			project := db.Project{
				ID:                fmt.Sprintf("%s%s", prefix, i.GuildID),
				Name:              name,
				Prefix:            prefix,
				RepoLink:          "",
				CategoryChannelID: category.ID,
				IssueChannelID:    issuesChannel.ID,
				AutoListMessageID: alMsg.ID,
				GuildID:           i.GuildID,
			}

			result := global.DB.Table("projects").Create(&project)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrDuplicatedKey) { // FIX:
					return ErrProjectAlreadyExists
				}
				return result.Error
			}

			return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
				Title:       fmt.Sprintf("Project %s created", name),
				Description: fmt.Sprintf("Check out <#%s>", issuesChannel.ID),
			}, false)
		case "rename":
			name := options["name"].StringValue()
			project.Name = name
			err := global.DB.Save(&project).Error
			if err != nil {
				return err
			}

			_, err = s.ChannelEdit(project.CategoryChannelID, &dg.ChannelEdit{Name: name})
			if err != nil {
				return err
			}

			return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
				Title: fmt.Sprintf("Project %s renamed", name),
			}, false)
		case "delete":
			confirmation := options["confirm"].BoolValue()

			if !confirmation {
				return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
					Title: fmt.Sprintf("alright, no actions taken."),
				}, true)
			}

			// delete in DB
			deleteMeProject := db.Project{}
			err := global.DB.Delete(&deleteMeProject, "prefix = ?", prefix).Error
			if err != nil {
				return err
			}

			deleteMeIssue := db.Issue{}
			err = global.DB.Delete(&deleteMeIssue, "project_id = ?", project.ID).Error
			if err != nil {
				return err
			}

			// delete discord channels
			_, err = s.ChannelDelete(project.IssueChannelID)
			if err != nil {
				return err
			}
			_, err = s.ChannelDelete(project.CategoryChannelID)
			if err != nil {
				return err
			}

			return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
				Title: fmt.Sprintf("Project %s deleted", project.Name),
			}, false)

		case "resetlistmsg":
			// TODO: Resend AutoList
		}

		return nil

	},
}
