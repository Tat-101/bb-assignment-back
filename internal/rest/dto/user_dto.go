package dto

import (
	"time"

	"github.com/tat-101/bb-assignment-back/domain"
)

type UserDTO struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func FromUserEntity(user *domain.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

func FromUserEntities(users []domain.User) []UserDTO {
	userDTOs := make([]UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = FromUserEntity(&user)
	}
	return userDTOs
}
