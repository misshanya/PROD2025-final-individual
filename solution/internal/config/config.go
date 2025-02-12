package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		log.Println("Trying to get variables from environment...")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("missing DATABASE_URL env variable")
	}
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		log.Println("missing SERVER_ADDRESS env variable, using default (127.0.0.1:8080)")
		serverAddress = "127.0.0.1:8080"
	}

	return &Config{
		DatabaseURL:   dbURL,
		ServerAddress: serverAddress,
	}
}
