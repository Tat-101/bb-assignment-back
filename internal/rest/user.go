package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tat-101/bb-assignment-back/domain"
)

//go:generate mockery --name UserService
type UserService interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUserByID(id string, updatedUser domain.User) error
	DeleteUserByID(id string) error

	AuthenticateUser(email, password string) (string, error)
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(r *gin.Engine, svc UserService) {
	handler := &UserHandler{
		Service: svc,
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", handler.GetUsers)
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// TODO:
