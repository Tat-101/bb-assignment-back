package config

import (
	"log"
	"os"

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
	Version       string
}

// LoadEnv loads env vars from .env
func LoadConfig() Config {
	// FIXME: wtf
	path := ".env"
	failCount := 0
	for {
		err := godotenv.Load(path)
		if err != nil && failCount > 10 {
			// log.Printf("Error loading .env file %v", err)
			log.Printf("Error loading .env file %v", err)
			os.Exit(-1)
		}
		if err == nil {
			// log.Println("path", path)
			break
		}
		path = "../" + path
		failCount++
	}

	return Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "bb-assignment"),
		DBPort:        getEnv("DB_PORT", "5432"),
		ServerAddress: getEnv("SERVER_ADDRESS", "3000"),
		JWTSecret:     getEnv("JWT_SECRET", "my_secret"),
		Version:       getEnv("API_VERSION", "v0"),
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
