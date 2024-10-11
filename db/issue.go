package db

type Issue struct {
	ID          string `gorm:"primarykey"`
	Code        string // #103
	Title       string
	Description string
	AssigneeIDs string // comma separated ids

	Roles []Role `gorm:"many2many:issue_roles;"`

	ProjectID string
}
