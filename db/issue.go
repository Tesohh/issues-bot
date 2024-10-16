package db

import (
	"fmt"
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
	AssigneeIDs string // comma separated ids
	IssueStatus uint8

	ThreadID  string `gorm:"column:thread_id"`
	MessageID string `gorm:"column:message_id"` // the message that contains the embed

	Roles []Role `gorm:"many2many:issue_roles;"`

	ProjectID string
}

func (i *Issue) PrettyID(projectPrefix string, count int) string {
	return fmt.Sprintf("#%s-%d", projectPrefix, count+1)
}

func (i *Issue) ThreadName() string {
	return fmt.Sprintf("%s %s %s", i.ID, IssueStatusIcons[i.IssueStatus], i.Title)
}
