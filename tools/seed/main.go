package main

import (
	"log"

	"github.com/tat-101/bb-assignment-back/config"
	"github.com/tat-101/bb-assignment-back/database"
	"github.com/tat-101/bb-assignment-back/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdminUser seeds the database with an admin user if it doesn't exist.
func SeedAdminUser(db *gorm.DB) {
	var user domain.User
	email := "admin@bb.com"

	// Check if the admin user already exists
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		log.Println("Admin user already exists, skipping seeding.")
		return
	}

	// Create a new admin user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	admin := domain.User{
		Name:     "Admin",
		Email:    email,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	log.Println("Admin user created successfully.")
}

func main() {
	cfg := config.LoadConfig()

	db := database.Initialize(cfg)

	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	SeedAdminUser(db)
}
