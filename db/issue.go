package db

import (
	"fmt"
	"issues/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

const (
	IssueStatusTodo     = 0
	IssueStatusDoing    = 1
	IssueStatusDone     = 2
	IssueStatusCanceled = 3
)

var IssueStatusIcons = [4]string{"🟩", "🟦", "🟪", "🟥"}
var IssueStatusColors = [4]int{0x7cb45c, 0x54acee, 0xa98ed6, 0xdd2e44}
var IssueStatusNames = [4]string{"todo", "doing", "done", "canceled"}

type Issue struct {
	ID          string `gorm:"primarykey"`
	Title       string
	Description string
	RecruiterID string `gorm:"column:recruiter_id"`
	AssigneeIDs string // comma separated ids
	IssueStatus uint8

	ThreadID  string `gorm:"column:thread_id"`
	MessageID string `gorm:"column:message_id"` // the message that contains the embed

	Roles []Role `gorm:"many2many:issue_roles;"` // note: should always be [kind, priority].

	ProjectID string
}

func (i *Issue) PrettyID(projectPrefix string, count int) string {
	return fmt.Sprintf("#%s-%d", projectPrefix, count+1)
}

func (i *Issue) ThreadName() string {
	return fmt.Sprintf("%s %s %s", i.ID, IssueStatusIcons[i.IssueStatus], i.Title)
}

// NOTE: assumes the first role is kind and the second is priority.
func (i *Issue) Embed() *dg.MessageEmbed {
	assigneeIDs := strings.Split(i.AssigneeIDs, ",")
	mentions := slash.MentionMany(assigneeIDs, "@", ", ")

	return &dg.MessageEmbed{
		Title: fmt.Sprintf("%s %s", i.ID, i.Title),
		Description: fmt.Sprintf(`
            **Kind**: <@&%s>
            **Priority**: <@&%s>
            **Recruiter**: <@%s>
            **Assignee(s)**: %s 

            %s
            `, i.Roles[0].ID, i.Roles[1].ID, i.RecruiterID, mentions, i.Description),
		Color: slash.EmbedColor,
	}
}
