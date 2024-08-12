package repository

import (
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/tools"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserByID(id string, updatedUser domain.User) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	user.Name = tools.Coalesce(updatedUser.Name, user.Name)
	if updatedUser.Password != "" {
		user.Password = updatedUser.Password
		if err := user.HashPassword(); err != nil {
			return nil, err
		}
	}

	return &user, r.DB.Save(&user).Error
}

func (r *UserRepository) DeleteUserByID(id string) error {
	return r.DB.Delete(&domain.User{}, id).Error
}
