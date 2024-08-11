package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint      `gorm:"primary_key"`
	Name      string    `gorm:"size:255"`
	Email     string    `gorm:"size:255;unique"`
	Password  string    `gorm:"size:255"`
	Role      string    `gorm:"size:50"` // "admin" or "user"
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
