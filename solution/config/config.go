package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerAddress string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	serverAddress := os.Getenv("SERVER_ADDRESS")

	return &Config{
		DatabaseURL: dbURL,
		ServerAddress: serverAddress,
	}
}