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
	Redis         RedisConfig
	OpenAI        OpenAIConfig
	MinIO         MinIOConfig
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type OpenAIConfig struct {
	BaseURL string
	ApiKey  string
}

type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicHost      string
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

	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	if openAIBaseURL == "" {
		log.Fatalln("missing OPENAI_BASE_URL")
	}

	openAIApiKey := os.Getenv("OPENAI_API_KEY")
	if openAIApiKey == "" {
		log.Println("[WARNING] OpenAI api key unset")
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		log.Fatalln("missing MINIO_ENDPOINT")
	}

	minioAccessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	if minioAccessKeyID == "" {
		log.Println("[WARNING] MinIO access key id unset")
	}

	minioSecretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	if minioSecretAccessKey == "" {
		log.Println("[WARNING] MinIO secret access key unset")
	}

	minioBucketName := os.Getenv("MINIO_BUCKET")
	if minioBucketName == "" {
		log.Println("[WARNING] MinIO bucket name unset, using default (proood)")
		minioBucketName = "proood"
	}

	minioPublicHost := os.Getenv("MINIO_PUB_HOST")
	if minioPublicHost == "" {
		log.Println("![WARNING]! MinIO public host unset, using default (localhost:9000)")
		minioPublicHost = "localhost:9000"
	}

	return &Config{
		DatabaseURL:   dbURL,
		ServerAddress: serverAddress,
		Redis: RedisConfig{
			Address:  redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
		OpenAI: OpenAIConfig{
			BaseURL: openAIBaseURL,
			ApiKey:  openAIApiKey,
		},
		MinIO: MinIOConfig{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKeyID,
			SecretAccessKey: minioSecretAccessKey,
			BucketName:      minioBucketName,
			PublicHost:      minioPublicHost,
		},
	}
}
