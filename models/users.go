package models

import (
	"errors"
	"time"

	sql "coderero.dev/projects/go/gin/hello/db"
	"gorm.io/gorm"
)

var db *gorm.DB

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

func init() {
	db = sql.GetDB()
}

func (u *User) Create() *User {
	db.Model(&u).Create(&u)
	return u
}

func (u User) CheckForUser(username string, email string) bool {
	var check User
	db.Model(&u).Where("username = ? OR email = ?", username, email).First(&check)
	if check.ID != 0 {
		return true
	}
	return false
}

func (u *User) GetUserForLogin(username, email, password string) *User {
	db.Model(&u).Where("username = ? OR email = ?", username, email).First(&u)
	return u
}

func (u *User) GetUserById(id int) *User {
	db.Model(&u).Where("id = ?", id).First(&u)
	return u
}

func (u *User) GetUserByEmail(email string) *User {
	db.Model(&u).Where("email = ?", email).First(&u)
	return u
}

func (u *User) GetUserByUsername(username string) *User {
	db.Model(&u).Where("username = ?", username).First(&u)
	return u
}

func (u *User) Update(id int) (*User, error) {
	if b, err := checkForId(id); b {
		return nil, err
	}

	db.Model(&u).Where("id = ?", id).Updates(&u)
	return u, nil
}

func (u *User) Delete(id int) error {
	if b, err := checkForId(id); b {
		return err
	}

	db.Model(&u).Where("id = ?", id).Delete(&u)
	return nil
}
