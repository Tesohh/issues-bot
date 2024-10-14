package handler

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

// TODO: split this function more
func AddIssue(s *dg.Session, m *dg.MessageCreate, roleIDs, channelIDs, userIDs []string, title, desc string) error {
	// get current channel/thread
	currentChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	// get project from channelid or parent channelid
	var project db.Project
	result := global.DB.First(&project, "issue_channel_id = ? or issue_channel_id = ?", m.ChannelID, currentChannel.ParentID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotInProject
		}
		return result.Error
	}

	// get the actual channel from the project
	ch, err := s.Channel(project.IssueChannelID)
	if err != nil {
		return err
	}

	// make the thread
	var count int64
	global.DB.Table("issues").Where("project_id = ?", project.ID).Count(&count)

	issueID := fmt.Sprintf("#%s-%d", project.Prefix, count+1)
	threadName := fmt.Sprintf("%s %s", issueID, title)
	thread, err := s.ThreadStart(ch.ID, threadName, dg.ChannelTypeGuildPublicThread, 10080)
	if err != nil {
		return err
	}

	// parse all roles and get the kind and priority
	kindRole, priorityRole, err := ParseRoles(s, m.GuildID, roleIDs)
	if err != nil {
		return err
	}
	assignees := ParseAssignees(m.Author.ID, userIDs)

	issue := db.Issue{
		ID:          issueID,
		Title:       title,
		Description: desc,
		AssigneeIDs: strings.Join(assignees, ","),
		ThreadID:    thread.ID,
		Roles:       []db.Role{*kindRole, *priorityRole},
		ProjectID:   project.ID,
	}

	result = global.DB.Table("issues").Create(&issue)
	if result.Error != nil {
		return result.Error
	}

	var roleMentions []string
	for _, id := range assignees {
		mention := fmt.Sprintf("<@%s>", id)
		roleMentions = append(roleMentions, mention)
	}

	embed := dg.MessageEmbed{
		Title: threadName,
		Description: fmt.Sprintf(`
            **Kind**: <@&%s>
            **Priority**: <@&%s>
            **Recruiter**: <@%s>
            **Assignee(s)**: %s 

            %s
            `, kindRole.ID, priorityRole.ID, m.Author.ID, strings.Join(roleMentions, ", "), desc),
		Color: slash.EmbedColor,
	}

	roleMentions = append(roleMentions, fmt.Sprintf("<@%s>", m.Author.ID))

	// temporarily mention users to make the thread show up
	deleteMeMsg, err := s.ChannelMessageSend(thread.ID, strings.Join(roleMentions, ""))
	if err != nil {
		return err
	}
	_ = s.ChannelMessageDelete(thread.ID, deleteMeMsg.ID)

	embedMsg, err := s.ChannelMessageSendEmbed(thread.ID, &embed)
	if err != nil {
		return err
	}

	err = s.ChannelMessagePin(thread.ID, embedMsg.ID)
	if err != nil {
		return err
	}

	// notify mentioned channels
	err = NotifyMentionedChannels(s, m, thread.ID, channelIDs)
	if err != nil {
		return err
	}

	return nil
}
