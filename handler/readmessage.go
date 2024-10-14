package handler

import (
	"log/slog"
	"regexp"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var roleMentionRegex = regexp.MustCompile(`<@&([0-9]+)>`)
var channelMentionRegex = regexp.MustCompile(`<#([0-9]+)>`)
var userMentionRegex = regexp.MustCompile(`<@([0-9]+)>`)

func ReadMessage(s *dg.Session, m *dg.MessageCreate) {
	str := m.Content
	if len(str) == 0 {
		return
	}
	if str[0] != '-' {
		return
	}
	str = strings.TrimLeft(str, "- ")

	roleMatches := roleMentionRegex.FindAllStringSubmatch(str, -1)
	roleIDs := make([]string, 0)
	for _, match := range roleMatches {
		roleIDs = append(roleIDs, match[1])
	}

	channelMatches := channelMentionRegex.FindAllStringSubmatch(str, -1)
	channelIDs := make([]string, 0)
	for _, match := range channelMatches {
		channelIDs = append(channelIDs, match[1])
	}

	userMatches := userMentionRegex.FindAllStringSubmatch(str, -1)
	userIDs := make([]string, 0)
	for _, match := range userMatches {
		userIDs = append(userIDs, match[1])
	}

	str = roleMentionRegex.ReplaceAllString(str, "")
	// str = channelMentionRegex.ReplaceAllString(str, "") # dont remove channels...
	str = userMentionRegex.ReplaceAllString(str, "")
	str = strings.Trim(str, " ")

	splits := strings.Split(str, "\n")
	title := splits[0]
	description := strings.Join(splits[1:], "\n")

	slog.Info("received message:", "roleIDs", roleIDs, "channelIDs", channelIDs, "userIDs", userIDs, "str", str)
	err := AddIssue(s, m, roleIDs, channelIDs, userIDs, title, description)
	if err != nil {
		slog.Error(err.Error())
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}
}
