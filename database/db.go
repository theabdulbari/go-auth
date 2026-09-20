package database

import (
	"go-auth/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	// Both models must be listed
	if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		panic("Failed to migrate: " + err.Error())
	}

	DB = db
}
