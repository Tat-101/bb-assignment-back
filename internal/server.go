package internal

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tat-101/bb-assignment-back/config"
	"github.com/tat-101/bb-assignment-back/database"
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/internal/repository"
	"github.com/tat-101/bb-assignment-back/internal/rest"
	"github.com/tat-101/bb-assignment-back/user"
)

func SetupServer() *gin.Engine {
	cfg := config.LoadConfig()

	db := database.Initialize(cfg)

	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	userRepo := repository.NewUserRepository(db)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost:") {
				return true
			}
			if origin == "https://bb-assignment-front.pages.dev" {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": cfg.Version,
		})
	})

	userService := user.NewService(userRepo)
	rest.NewUserHandler(r, userService)

	return r
}
