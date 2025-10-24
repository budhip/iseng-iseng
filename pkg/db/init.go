package db

import (
	// "os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// PLEASE USE SQLITE, WE USE GORM TO MAKE IT EASIER FOR YOU, BUT YOU CAN CHANGE IT
// InitAndMigrate initializes the database connection and creates required tables.
func InitAndMigrate(models ...interface{}) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error
	// Remove existing DB file if it exists
	//_ = os.Remove("amartha.db")

	db, err = gorm.Open(sqlite.Open("amartha.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, err
	}

	return db, nil
}
