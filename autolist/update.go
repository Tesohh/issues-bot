package autolist

import (
	"fmt"
	"issues/db"
	"issues/global"

	"github.com/bwmarrin/discordgo"
)

// if project || guild is nil, it will request it from the db.
func Update(project *db.Project, projectId string, guild *db.Guild, guildID string, s *discordgo.Session, lastUpdate string) error {
	if project == nil {
		err := global.DB.
			Preload("Issues").
			First(&project, "id = ?", projectId).Error
		if err != nil {
			return err
		}
	}
	if guild == nil {
		err := global.DB.
			First(&guild, "id = ?", guildID).Error
		if err != nil {
			return err
		}
	}
	embedTitle := fmt.Sprintf("AutoList™️ for %s", project.Name)
	filteredIssues := ApplyFilters(project.Issues, false, "", false, "", "")

	autolistEmbed := Embed(embedTitle, guild.DefaultPriorityRoleID, filteredIssues, lastUpdate)
	_, err := s.ChannelMessageEditEmbed(project.IssueChannelID, project.AutoListMessageID, &autolistEmbed)
	if err != nil {
		return err
	}

	return nil
}
