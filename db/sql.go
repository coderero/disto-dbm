package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// `var db *gorm.DB` is declaring a variable named `db` of type `*gorm.DB`. This variable will hold a
// pointer to a `gorm.DB` object.
var db *gorm.DB

// The `init` function initializes a connection to a SQLite database using the GORM library in Go.
func init() {
	godotenv.Load()
	uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db_local, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db = db_local
}

// The GetDB function returns a pointer to a gorm.DB object.
func GetDB() *gorm.DB {
	return db
}
