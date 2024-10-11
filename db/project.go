package db

type Project struct {
	ID       string `gorm:"primarykey"`
	Name     string
	RepoLink string

	GuildID string

	Issues []Issue
}
