package handler

import (
	"fmt"
	"issues/db"
	"issues/global"

	dg "github.com/bwmarrin/discordgo"
)

func NotifyMentionedChannels(s *dg.Session, m *dg.MessageCreate, threadID string, channelIDs []string) error {
	issues := []db.Issue{}
	for _, channelID := range channelIDs {
		var issue db.Issue
		result := global.DB.Where("thread_id = ?", channelID).Find(&issue)
		if result.Error != nil {
			return result.Error
		}
		issues = append(issues, issue)
	}

	for _, issue := range issues {
		_, err := s.ChannelMessageSend(issue.ThreadID, fmt.Sprintf("This issue was mentioned in <#%s>", threadID))
		if err != nil {
			return err
		}
	}

	return nil
}
