package db

import (
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	slog.Info("Connecting to db...")

	db, err := gorm.Open(sqlite.Open(".data/issues.db"), &gorm.Config{})
	if err != nil {
		return db, err
	}

	// Migrate the schemas
	db.AutoMigrate(&Guild{}, &Role{}, &Project{}, &Issue{})

	return db, nil
}
