package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port  string
	DbUrl string
}

func LoadConfig() *Config {
	return &Config{
		Port:  getEnv("PORT", "8080"),
		DbUrl: getDbUrl("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getDbUrl(key, fallback string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
