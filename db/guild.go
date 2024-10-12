package db

type Guild struct {
	ID         string `gorm:"primarykey"`
	Registered int64  `gorm:"autoCreateTime"`

	DefaultKindRoleID     string
	DefaultKindRole       Role `gorm:"foreignKey:DefaultKindRoleID"`
	DefaultPriorityRoleID string
	DefaultPriorityRole   Role `gorm:"foreignKey:DefaultPriorityRoleID"`

	Roles   []Role
	Project []Project
}
