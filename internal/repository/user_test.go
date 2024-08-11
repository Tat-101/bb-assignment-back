package repository_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tat-101/bb-assignment-back/config"
	"github.com/tat-101/bb-assignment-back/database"
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/internal/repository"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	cfg := config.LoadConfig()
	db := database.Initialize(cfg)

	err := db.AutoMigrate(&domain.User{})
	require.NoError(t, err)

	tx := db.Begin()

	return tx, func() {
		tx.Rollback()
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	userRepo := repository.NewUserRepository(db)

	user := &domain.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	err := userRepo.CreateUser(user)

	assert.NoError(t, err)

	var dbUser domain.User
	err = db.First(&dbUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, dbUser.Email)
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	userRepo := repository.NewUserRepository(db)

	users := []domain.User{
		{Email: "user1@example.com", Name: "User One"},
		{Email: "user2@example.com", Name: "User Two"},
	}

	for _, user := range users {
		err := userRepo.CreateUser(&user)
		require.NoError(t, err)
	}

	dbUsers, err := userRepo.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, dbUsers, len(users))
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	userRepo := repository.NewUserRepository(db)

	user := domain.User{Email: "test@example.com", Name: "Test User"}
	err := userRepo.CreateUser(&user)
	require.NoError(t, err)

	dbUser, err := userRepo.GetUserByID(user.ID)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, dbUser.Email)
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	userRepo := repository.NewUserRepository(db)

	user := domain.User{Email: "test@example.com", Name: "Test User"}
	err := userRepo.CreateUser(&user)
	require.NoError(t, err)

	dbUser, err := userRepo.GetUserByEmail(user.Email)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, dbUser.Email)
}

func TestUserRepository_UpdateUserByID(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	userRepo := repository.NewUserRepository(db)

	user := domain.User{Email: "test@example.com", Name: "Test User", Password: "password123"}
	err := userRepo.CreateUser(&user)
	require.NoError(t, err)

	updatedUser := domain.User{Name: "Updated User", Email: "updated@example.com", Password: "newpassword123"}

	err = userRepo.UpdateUserByID(fmt.Sprint(user.ID), updatedUser)
	assert.NoError(t, err)

	dbUser, err := userRepo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Email, dbUser.Email)
	assert.Equal(t, updatedUser.Name, dbUser.Name)
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	userRepo := repository.NewUserRepository(db)

	user := domain.User{Email: "test@example.com", Name: "Test User"}
	err := userRepo.CreateUser(&user)
	require.NoError(t, err)

	err = userRepo.DeleteUserByID(fmt.Sprint(user.ID))
	assert.NoError(t, err)

	_, err = userRepo.GetUserByID(user.ID)
	assert.Error(t, err) // Should return an error because the user has been deleted
}
