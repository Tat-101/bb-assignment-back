package main

import (
	"github.com/tat-101/bb-assignment-back/config"
	"github.com/tat-101/bb-assignment-back/database"
	"github.com/tat-101/bb-assignment-back/domain"
)

func main() {
	cfg := config.LoadConfig()

	db := database.Initialize(cfg)

	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

}
