package user

import (
	"errors"

	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/tools"
	"golang.org/x/crypto/bcrypt"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//go:generate mockery --name UserRepository
type UserRepository interface {
	CreateUser(user *domain.User) error
	// TODO: improve should have skip limit
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUserByID(id string, updatedUser domain.User) error
	DeleteUserByID(id string) error

	AuthenticateUser(email, password string) (*domain.User, error)
}

type Service struct {
	userRepo UserRepository
}

func NewService(u UserRepository) *Service {
	return &Service{
		userRepo: u,
	}
}

func (s *Service) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.GetAllUsers()
}

// CreateUser creates a new user in the repository
func (s *Service) CreateUser(user *domain.User) error {
	return s.userRepo.CreateUser(user)
}

// GetUserByID retrieves a user by their ID from the repository
func (s *Service) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.GetUserByID(id)
}

// GetUserByEmail retrieves a user by their email from the repository
func (s *Service) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

// UpdateUserByID updates a user's information by their ID in the repository
func (s *Service) UpdateUserByID(id string, updatedUser domain.User) error {
	return s.userRepo.UpdateUserByID(id, updatedUser)
}

// DeleteUserByID deletes a user by their ID from the repository
func (s *Service) DeleteUserByID(id string) error {
	return s.userRepo.DeleteUserByID(id)
}

func (s *Service) AuthenticateUser(loginData LoginData) (string, error) {
	user, err := s.userRepo.GetUserByEmail(loginData.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := tools.GenerateJWT(user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
