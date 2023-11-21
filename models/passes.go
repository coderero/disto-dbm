package models

type UsedPassword struct {
	ID       uint   `json:"-" gorm:"primarykey"`
	Password string `json:"-" gorm:"not null"`
	UserID   uint   `json:"-" gorm:"not null,user_id"`
}

// GetUsedPasswords returns all the used passwords for a user.
func (u *User) GetUsedPasswords() []UsedPassword {
	var usedPasswords []UsedPassword
	db.Model(&u).Where("user_id = ?", u.ID).Find(&usedPasswords)
	return usedPasswords
}

// AddUsedPassword adds a new used password for a user.
func (u *User) AddUsedPassword(password string) {
	db.Model(&u).Create(&UsedPassword{
		Password: password,
		UserID:   u.ID,
	})
}
