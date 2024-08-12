package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/internal/rest/dto"
	"github.com/tat-101/bb-assignment-back/internal/rest/middleware"
	"github.com/tat-101/bb-assignment-back/internal/rest/service"
)

type UserHandler struct {
	Service service.UserService
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(r *gin.Engine, svc service.UserService) {
	handler := &UserHandler{
		Service: svc,
	}

	authMiddleware := middleware.AuthMiddleware(svc)
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", authMiddleware, handler.GetUsers)
		userRoutes.POST("", authMiddleware, handler.CreateUser)
		userRoutes.GET("/:id", authMiddleware, handler.GetUserByID)
		userRoutes.PUT("/:id", authMiddleware, handler.UpdateUserByID)
		userRoutes.DELETE("/:id", authMiddleware, middleware.AdminMiddleware(), handler.DeleteUserByID)
	}

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", handler.LoginUser)
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.FromUserEntities(users))
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.FromUserEntity(&user))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	idStr, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.GetUserByID(uint(idStr))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, dto.FromUserEntity(user))
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	id := c.Param("id")
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.Service.UpdateUserByID(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.FromUserEntity(updatedUser))
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteUserByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginData LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Service.AuthenticateUser(loginData.Email, loginData.Password)
	if err != nil {
		// fmt.Println("22222222")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
