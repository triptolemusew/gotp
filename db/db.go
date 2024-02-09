package db

import (
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetClient(appDir string) (*gorm.DB, error) {
	homeDir, _ := os.UserHomeDir()
	dbPath := path.Join(homeDir, appDir, "gotp.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(appDir string) error {
	db, err := GetClient(appDir)
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&Key{}); err != nil {
		return err
	}

	return nil
}
