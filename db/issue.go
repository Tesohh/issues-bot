package db

type Issue struct {
	ID          string `gorm:"primarykey"`
	Title       string
	Description string
	AssigneeIDs string // comma separated ids

	ThreadID string

	Roles []Role `gorm:"many2many:issue_roles;"`

	ProjectID string
}
