package service

import "github.com/tat-101/bb-assignment-back/domain"

//go:generate mockery --name UserService
type UserService interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUserByID(id string, updatedUser domain.User) (*domain.User, error)
	DeleteUserByID(id string) error

	AuthenticateUser(email, password string) (string, error)
	ValidateToken(token string) (*domain.User, error)
}
