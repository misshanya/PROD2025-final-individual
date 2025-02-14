package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
	Redis RedisConfig
}

type RedisConfig struct {
	Address string
	Password string
	DB int
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

	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		log.Fatalln("missing REDIS_ADDRESS env variable")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		log.Println("[WARNING] redis password unset")
	}

	var redisDB int
	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr == "" {
		log.Println("missing REDIS_DB env variable, using default (0)")
	} else {
		redisDB, err = strconv.Atoi(redisDBStr)
		if err != nil {
			log.Fatalln("failed to convert redis db to integer")
		}
	}

	return &Config{
		DatabaseURL:   dbURL,
		ServerAddress: serverAddress,
		Redis: RedisConfig{
			Address: redisAddr,
			Password: redisPassword,
			DB: redisDB,
		},
	}
}
