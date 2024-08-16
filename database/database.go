package database

import (
	"log"

	"github.com/tat-101/bb-assignment-back/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Initialize(cfg config.Config) *gorm.DB {
	config.LoadConfig()

	dsn := cfg.GetDBConfig()

	// fmt.Println("dsn", dsn)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db = database

	return db
}

func GetDB() *gorm.DB {
	return db
}
