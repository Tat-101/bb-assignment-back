package main

import (
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

// TODO: improve log
func main() {
	cfg := config.LoadConfig()

	db := database.Initialize(cfg)

	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	userRepo := repository.NewUserRepository(db)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:*", "https://bb-assignment-front.pages.dev/*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           12 * time.Hour,
	}))

	userService := user.NewService(userRepo)
	rest.NewUserHandler(r, userService)

	r.Run(":" + cfg.ServerAddress)
}
