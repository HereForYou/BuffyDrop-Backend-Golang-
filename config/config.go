package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
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
