package user_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/tools"
	"github.com/tat-101/bb-assignment-back/user"
	"github.com/tat-101/bb-assignment-back/user/mocks"
	"golang.org/x/crypto/bcrypt"
)

func TestService_GetAllUsers(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	expectedUsers := []domain.User{
		{ID: 1, Email: "user1@example.com", Name: "User One"},
		{ID: 2, Email: "user2@example.com", Name: "User Two"},
	}

	mockUserRepo.On("GetAllUsers").Return(expectedUsers, nil)

	users, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockUserRepo.AssertExpectations(t)
}

func TestService_CreateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	newUser := &domain.User{Email: "newuser@example.com", Name: "New User"}

	mockUserRepo.On("CreateUser", newUser).Return(nil)

	err := service.CreateUser(newUser)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestService_GetUserByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	expectedUser := &domain.User{ID: 1, Email: "user@example.com", Name: "User One"}

	mockUserRepo.On("GetUserByID", uint(1)).Return(expectedUser, nil)

	user, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertExpectations(t)
}

func TestService_GetUserByEmail(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	expectedUser := &domain.User{Email: "user@example.com", Name: "User One"}

	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(expectedUser, nil)

	user, err := service.GetUserByEmail("user@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertExpectations(t)
}

func TestService_UpdateUserByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	updatedUser := domain.User{Email: "updated@example.com", Name: "Updated User"}

	mockUserRepo.On("UpdateUserByID", "1", updatedUser).Return(&updatedUser, nil)

	newUser, err := service.UpdateUserByID("1", updatedUser)

	assert.NoError(t, err)
	assert.Equal(t, newUser.Name, updatedUser.Name)
	mockUserRepo.AssertExpectations(t)
}

func TestService_DeleteUserByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	mockUserRepo.On("DeleteUserByID", "1").Return(nil)

	err := service.DeleteUserByID("1")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestService_AuthenticateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	email := "user"
	password := "password123"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	expectedUser := &domain.User{
		Email:    "user@example.com",
		Password: string(hashedPassword),
	}

	mockUserRepo.On("GetUserByEmail", email).Return(expectedUser, nil)

	token, err := service.AuthenticateUser(email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestService_AuthenticateUser_Fail(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	email := "user"
	password := "wrongpassword"

	expectedUser := &domain.User{
		Email:    "user@example.com",
		Password: "$2a$10$12345678901234567890123456789012345678901234567890",
	}

	mockUserRepo.On("GetUserByEmail", email).Return(expectedUser, nil)

	token, err := service.AuthenticateUser(email, password)

	assert.Error(t, err)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestService_ValidateToken_Success(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	email := "test@example.com"
	token, err := tools.GenerateJWT(email)
	// fmt.Println("token", token)

	assert.NoError(t, err)

	expectedUser := &domain.User{
		Email: email,
		Name:  "Test User",
	}

	mockUserRepo.On("GetUserByEmail", email).Return(expectedUser, nil)

	user, err := service.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertExpectations(t)
}

func TestService_ValidateToken_InvalidToken(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	user, err := service.ValidateToken("invalidToken")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid token", err.Error())
}

func TestService_ValidateToken_UserNotFound(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	service := user.NewService(mockUserRepo)

	email := "notfound@example.com"
	token, _ := tools.GenerateJWT(email)

	mockUserRepo.On("GetUserByEmail", email).Return(nil, errors.New("user not found"))

	user, err := service.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}
