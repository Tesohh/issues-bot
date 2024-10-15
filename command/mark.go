package command

import (
	"fmt"
	"issues/db"
	"issues/global"
	"issues/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var Mark = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "mark",
		Description: "mark issue as ...",
		Options: []*dg.ApplicationCommandOption{
			{Type: dg.ApplicationCommandOptionSubCommand, Name: "todo", Description: "� todo"},
			{Type: dg.ApplicationCommandOptionSubCommand, Name: "doing", Description: "� doing"},
			{Type: dg.ApplicationCommandOptionSubCommand, Name: "done", Description: "� done"},
			{Type: dg.ApplicationCommandOptionSubCommand, Name: "cancelled", Description: "� cancelled"},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0].Name

		var issueStatus uint8
		var archive = false
		var notifyAssignees = false
		var lock = false
		switch subcommand {
		case "todo":
			issueStatus = db.IssueStatusTodo
		case "doing":
			issueStatus = db.IssueStatusDoing
		case "done":
			issueStatus = db.IssueStatusDone
			archive = true
			notifyAssignees = true
		case "cancelled":
			issueStatus = db.IssueStatusCanceled
			archive = true
			notifyAssignees = true
			lock = true
		}

		// get issue from current channel id
		var issue db.Issue
		result := global.DB.Find(&issue, "thread_id = ?", i.ChannelID)
		if result.Error != nil {
			return result.Error
		}

		if issue.ID == "" {
			return ErrNotInIssue
		}

		if issue.IssueStatus == issueStatus {
			embed := dg.MessageEmbed{
				Title: "status is unchanged, doing nothing",
			}
			return slash.ReplyWithEmbed(s, i, embed, true)
		}

		// modify issue
		issue.IssueStatus = issueStatus
		err := global.DB.Save(&issue).Error
		if err != nil {
			return err
		}

		// change thread name
		_, err = s.ChannelEdit(issue.ThreadID, &dg.ChannelEdit{
			Name:     issue.ThreadName(),
			Archived: &archive,
			Locked:   &lock,
		})
		if err != nil {
			return err
		}

		// notify assignees if needed
		if notifyAssignees {
			tagsString := slash.MentionMany(strings.Split(issue.AssigneeIDs, ","), "@", " ")
			deleteMeMsg, err := s.ChannelMessageSend(i.ChannelID, tagsString)
			if err != nil {
				return err
			}

			_ = s.ChannelMessageDelete(i.ChannelID, deleteMeMsg.ID) // dont really care if it cant delete it
		}

		// and finally send the embed
		embed := dg.MessageEmbed{
			Title: fmt.Sprintf("Marked as %s", db.IssueStatusNames[issueStatus]),
			Color: db.IssueStatusColors[issueStatus],
		}

		return slash.ReplyWithEmbed(s, i, embed, false)
	},
}
