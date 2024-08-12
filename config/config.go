package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	ServerAddress string
	JWTSecret     string
}

// LoadEnv loads env vars from .env
func LoadConfig() Config {
	re := regexp.MustCompile(`^(.*` + "bb-assignment-back" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {

		log.Fatalf("Error loading .env file %v", err)
		os.Exit(-1)
	}

	return Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "mydb"),
		DBPort:        getEnv("DB_PORT", "5432"),
		ServerAddress: getEnv("SERVER_ADDRESS", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
	}
}

// getEnv fetches the value of an environment variable, or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (cfg Config) GetDBConfig() string {
	return "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
}
