package command

import (
	"errors"
	"fmt"
	"issues/db"
	"issues/global"
	"issues/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var NewProject = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "newproject",
		Description: "creates a new project for you",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionString,
				Name:        "name",
				Description: "the project's name",
				Required:    true,
			},
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
	Func: func(s *dg.Session, i *dg.Interaction) error {
		options := slash.GetOptionMap(i)
		name := options["name"].StringValue()
		prefix := strings.ToUpper(options["prefix"].StringValue())

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

		project := db.Project{
			ID:                fmt.Sprintf("%s%s", prefix, i.GuildID),
			Name:              name,
			Prefix:            prefix,
			RepoLink:          "",
			CategoryChannelID: category.ID,
			IssueChannelID:    issuesChannel.ID,
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
		})
	},
}
