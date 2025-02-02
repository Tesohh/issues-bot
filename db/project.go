package db

type Project struct {
	ID       string `gorm:"primarykey"`
	Name     string
	Prefix   string
	RepoLink string

	CategoryChannelID string
	IssueChannelID    string `gorm:"size:3"`
	AutoListMessageID string
	GuildID           string

	Issues []Issue
}
