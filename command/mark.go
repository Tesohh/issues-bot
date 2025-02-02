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
		var autoArchiveDuration = 10080
		switch subcommand {
		case "todo":
			issueStatus = db.IssueStatusTodo
		case "doing":
			issueStatus = db.IssueStatusDoing
		case "done":
			issueStatus = db.IssueStatusDone
			archive = true
			notifyAssignees = true
			autoArchiveDuration = 60
		case "cancelled":
			issueStatus = db.IssueStatusCanceled
			archive = true
			notifyAssignees = true
			lock = true
			autoArchiveDuration = 60
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
			Name:                issue.ThreadName(),
			AutoArchiveDuration: autoArchiveDuration,
			Archived:            &archive,
			Locked:              &lock,
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
		var alsoWillArchiveString string
		if archive {
			alsoWillArchiveString = ", thread will be archived in 1 hour if inactive"
		}

		embed := dg.MessageEmbed{
			Title: fmt.Sprintf("Marked as %s%s", db.IssueStatusNames[issueStatus], alsoWillArchiveString),
			Color: db.IssueStatusColors[issueStatus],
		}

		err = slash.ReplyWithEmbed(s, i, embed, false)
		if err != nil {
			return err
		}

		err = autolist.Update(nil, issue.ProjectID, nil, i.GuildID, s, fmt.Sprintf("marked %s as %s", issue.ID, db.IssueStatusNames[issueStatus]))
		if err != nil {
			return err
		}

		return nil
	},
}
