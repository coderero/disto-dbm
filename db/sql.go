package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// `var db *gorm.DB` is declaring a variable named `db` of type `*gorm.DB`. This variable will hold a
// pointer to a `gorm.DB` object.
var (
	db    *gorm.DB
	onece = sync.Once{}
)

// The `init` function initializes a connection to a SQLite database using the GORM library in Go.
func init() {
	// The `godotenv.Load()` function is used to load the environment variables from the `.env` file.
	godotenv.Load()

	// The `sync.Once` type is used to perform initialization tasks only once. The `sync.Once` type
	// provides a `Do()` method that takes a function as an argument and executes it only once.
	onece.Do(func() {
		// The `uri` variable is used to store the connection string for the database.
		uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

		// The `gorm.Open()` function is used to open a connection to the database. It takes the name of
		// the database driver and the connection string as arguments. It returns a pointer to a `gorm.DB`
		// object and an error.
		db_local, err := gorm.Open(postgres.Open(uri), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Silent),
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			panic("failed to connect database")
		}
		db = db_local
	})
}

// The GetDB function returns a pointer to a gorm.DB object.
func GetDB() *gorm.DB {
	return db
}
