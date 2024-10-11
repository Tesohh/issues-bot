package db

type Guild struct {
	ID         string `gorm:"primarykey"`
	Registered int64  `gorm:"autoCreateTime"`

	Roles []Role
}
