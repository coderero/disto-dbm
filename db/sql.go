package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// `var db *gorm.DB` is declaring a variable named `db` of type `*gorm.DB`. This variable will hold a
// pointer to a `gorm.DB` object.
var db *gorm.DB

// The `init` function initializes a connection to a SQLite database using the GORM library in Go.
func init() {
	db_local, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db_local
}

// The GetDB function returns a pointer to a gorm.DB object.
func GetDB() *gorm.DB {
	return db
}
