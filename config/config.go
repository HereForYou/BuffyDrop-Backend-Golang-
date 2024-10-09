package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port     string
	DbUrl    string
	BotToken string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var port string
	var dbUrl string
	var botToken string
	var exists bool

	if port, exists = os.LookupEnv("PORT"); !exists {
		port = "8080"
	}

	if dbUrl, exists = os.LookupEnv("DB_URL"); !exists {
		log.Fatal("No database Url in .env file")
	}

	if botToken, exists = os.LookupEnv("TG_BOT_TOKEN"); !exists {
		log.Fatal("No TG bot token in .env file")
	}

	return &Config{
		Port:     port,
		DbUrl:    dbUrl,
		BotToken: botToken,
	}
}
