package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"size:255;not null" binding:"required"`
	Email     string `gorm:"size:255;unique" binding:"required"`
	Password  string `gorm:"size:255;not null"`
	Role      string `gorm:"size:50;default:user"` // "admin" or "user"
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
