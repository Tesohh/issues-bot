package db

type Project struct {
	ID       string `gorm:"primarykey"`
	Name     string
	Prefix   string
	RepoLink string

	IssueChannelID string
	GuildID        string

	Issues []Issue
}
