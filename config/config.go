package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// LoadEnv loads env vars from .env
func LoadConfig() {
	re := regexp.MustCompile(`^(.*` + "bb-assignment-back" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {

		log.Fatalf("Error loading .env file %v", err)
		os.Exit(-1)
	}
}

func GetDBConfig() string {
	return os.Getenv("DB_CONNECTION_STRING")
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
