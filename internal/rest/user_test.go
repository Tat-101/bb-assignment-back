package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/internal/rest"
	"github.com/tat-101/bb-assignment-back/internal/rest/service/mocks"
)

func TestHandler_GetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// var mockUser domain.User
	// err := faker.FakeData(&mockUser)
	// fmt.Printf("%+v\n", mockUsers)
	// assert.NoError(t, err)
	mockUserService := new(mocks.UserService)
	userHandler := rest.UserHandler{Service: mockUserService}

	// mockListUser := make([]domain.User, 0)
	// mockListUser = append(mockListUser, mockUser)
	mockListUser := []domain.User{
		{
			ID:        10,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Date(2040, 7, 10, 0, 38, 44, 0, time.FixedZone("UTC+7", 7*3600)),
		},
		{
			ID:        20,
			Name:      "Jane Doe",
			Email:     "jane@example.com",
			CreatedAt: time.Date(2040, 7, 10, 0, 39, 44, 0, time.FixedZone("UTC+7", 7*3600)),
		},
	}

	mockUserService.On("GetAllUsers").Return(mockListUser, nil)

	router := gin.Default()
	router.GET("/users", userHandler.GetUsers)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	// fmt.Println(w.Body.String())
	expectedResponse := "["
	for i, user := range mockListUser {
		if i > 0 {
			expectedResponse += ","
		}
		createdAt := user.CreatedAt.Format(time.RFC3339)
		expectedResponse += `{"id":` + strconv.Itoa(int(user.ID)) + `,"email":"` + user.Email + `","name":"` + user.Name + `","createdAt":"` + createdAt + `"}`
	}
	expectedResponse += "]"
	// expectedResponse := `[{"id":1,"email":"john@example.com","name":"John Doe"},{"id":2,"email":"jane@example.com","name":"Jane Doe"}]`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	// Ensure the mock service was called as expected
	mockUserService.AssertExpectations(t)
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.UserService)
	userHandler := rest.UserHandler{Service: mockUserService}

	mockUser := domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockUserService.On("CreateUser", &mockUser).Return(nil)

	router := gin.Default()
	router.POST("/users", userHandler.CreateUser)

	body := `{"name":"John Doe", "email":"john@example.com", "password":"password123"}`
	req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	expectedResponse := `{"id":0,"email":"john@example.com","name":"John Doe","createdAt":"` + mockUser.CreatedAt.Format(time.RFC3339) + `"}`
	// fmt.Println(w.Body.String())
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockUserService.AssertExpectations(t)
}

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.UserService)
	userHandler := rest.UserHandler{Service: mockUserService}

	mockUser := domain.User{
		ID:        10,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Date(2040, 7, 10, 0, 38, 44, 0, time.FixedZone("UTC+7", 7*3600)),
	}

	mockUserService.On("GetUserByID", uint(10)).Return(&mockUser, nil)

	router := gin.Default()
	router.GET("/users/:id", userHandler.GetUserByID)

	req, _ := http.NewRequest(http.MethodGet, "/users/10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{"id":10,"email":"john@example.com","name":"John Doe","createdAt":"` + mockUser.CreatedAt.Format(time.RFC3339) + `"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockUserService.AssertExpectations(t)
}

func TestUserHandler_UpdateUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.UserService)
	userHandler := rest.UserHandler{Service: mockUserService}

	mockUser := domain.User{
		ID:    10,
		Name:  "John Smith",
		Email: "john@example.com",
	}

	mockUserService.On("UpdateUserByID", "10", mockUser).Return(&mockUser, nil)

	router := gin.Default()
	router.PUT("/users/:id", userHandler.UpdateUserByID)

	body := `{"name":"John Smith", "email":"john@example.com"}`
	req, _ := http.NewRequest(http.MethodPut, "/users/10", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// fmt.Println("asdfasfasdf", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{"id":10,"email":"john@example.com","name":"John Smith","createdAt":"` + mockUser.CreatedAt.Format(time.RFC3339) + `"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockUserService.AssertExpectations(t)
}

func TestUserHandler_DeleteUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.UserService)
	userHandler := rest.UserHandler{Service: mockUserService}

	mockUserService.On("DeleteUserByID", "10").Return(nil)

	router := gin.Default()
	router.DELETE("/users/:id", userHandler.DeleteUserByID)

	req, _ := http.NewRequest(http.MethodDelete, "/users/10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockUserService.AssertExpectations(t)
}

func TestUserHandler_LoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.UserService)
	userHandler := rest.UserHandler{Service: mockUserService}

	mockToken := "mockToken123"
	mockUserService.On("AuthenticateUser", "john@example.com", "password123").Return(mockToken, nil)

	router := gin.Default()
	router.POST("/auth/login", userHandler.LoginUser)

	body := `{"email":"john@example.com", "password":"password123"}`
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{"token":"mockToken123"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockUserService.AssertExpectations(t)
}
