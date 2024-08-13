package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"size:255;not null" faker:"name"`
	Email     string `gorm:"size:255;unique" faker:"email"`
	Password  string `gorm:"size:255;not null" faker:"password"`
	Role      string `gorm:"size:50;default:user"` // "admin" or "user"
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TODO: validation request, tag binding:"required"

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
