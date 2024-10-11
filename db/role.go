package db

type RoleType string

const (
	RoleTypeKind     RoleType = "kind"
	RoleTypePriority RoleType = "priority"
)

type Role struct {
	ID       string `gorm:"primarykey"`
	RoleType RoleType

	GuildId string
}
