package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tat-101/bb-assignment-back/internal/rest/middleware"
	"github.com/tat-101/bb-assignment-back/internal/rest/service"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(r *gin.Engine, svc service.UserService) {
	handler := &UserHandler{
		Service: svc,
	}

	authMiddleware := middleware.AuthMiddleware(svc)
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", handler.GetUsers)

		userRoutes.DELETE("/:id", authMiddleware, handler.GetUsers)
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
