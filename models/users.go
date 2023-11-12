package models

import (
	"errors"
	"time"

	sql "coderero.dev/projects/go/gin/hello/db"
	"gorm.io/gorm"
)

// The `init` function initializes the database connection and performs automatic migration for the
// `User` model.
var db *gorm.DB

func init() {
	db = sql.GetDB()
	db.AutoMigrate(&User{})
}

// The User struct defines the structure of a user record in the database.
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username,omitempty" gorm:"unique;not null"`
	Email     string         `json:"email,omitempty" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"firstname,omitempty" gorm:"not null"`
	LastName  string         `json:"lastname,omitempty" gorm:"not null"`
	Age       int            `json:"age,omitempty" gorm:"not null"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// The above code defines a User struct and provides methods for creating, retrieving, updating, and
// deleting user records in a database.
func checkForId(id int) (bool, error) {
	if id == 0 {
		return false, nil
	}
	return true, errors.New("'Id' should not be passed through the struct body instead use function parameter")
}

// The `Create()` method is used to create a new user record in the database. It takes a pointer to a
// `User` struct as its receiver (`u *User`) and returns a pointer to the created user (`*User`).
func (u *User) Create() *User {
	db.Model(&u).Create(&u)
	return u
}

// The `CheckForUser` method is used to check if a user with the given username or email already
// exists in the database.
func (u User) CheckForUser(username string, email string) bool {
	var check User
	db.Model(&u).Where("username = ? OR email = ?", username, email).First(&check)
	if check.ID != 0 {
		return true
	}
	return false
}

// The `GetUserForLogin` method is used to retrieve a user record from the database based on the
// provided username or email. It takes the username and email as parameters and returns a pointer to
// the retrieved user (`*User`).
func (u *User) GetUserForLogin(username, email string) *User {
	db.Model(&u).Raw("SELECT * FROM users WHERE username = ? OR email = ?", username, email).First(&u)
	return u
}

// The `GetUserById` method is used to retrieve a user record from the database based on the provided
// user ID. It takes the user ID as a parameter and returns a pointer to the retrieved user (`*User`).
func (u *User) GetUserById(id int) *User {
	db.Model(&u).Where("id = ?", id).First(&u)
	return u
}

// The `GetUserByEmail` method is used to retrieve a user record from the database based on the
// provided email. It takes the email as a parameter and returns a pointer to the retrieved user
// (`*User`).
func (u *User) GetUserByEmail(email string) *User {
	db.Model(&u).Where("email = ?", email).First(&u)
	return u
}

// The `GetUserByUsername` method is used to retrieve a user record from the database based on the
// provided username. It takes the username as a parameter and returns a pointer to the retrieved user
// (`*User`).
func (u *User) GetUserByUsername(username string) *User {
	db.Model(&u).Where("username = ?", username).First(&u)
	return u
}

// The `Update` method is a method defined on the `User` struct. It is used to update a user record in
// the database based on the provided user ID.
func (u *User) Update(id int) (*User, error) {
	if b, err := checkForId(id); b {
		return nil, err
	}

	db.Model(&u).Where("id = ?", id).Updates(&u)
	return u, nil
}

// The `Delete` method is a method defined on the `User` struct. It is used to delete a user record
// from the database based on the provided user ID.
func (u *User) Delete(id int) error {
	if b, err := checkForId(id); b {
		return err
	}

	db.Model(&u).Where("id = ?", id).Delete(&u)
	return nil
}
