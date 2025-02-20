package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/app"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/config"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/handlers"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/ml"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type Server struct {
	httpServer *http.Server
	db         *pgxpool.Pool
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	// Init db connection
	conn, err := InitDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Init SQL queries
	queries := storage.New(conn)

	// Init redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %v", err)
	}

	// Init MinIO client
	minioClient, err := initMinIO(
		ctx,
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKeyID,
		cfg.MinIO.SecretAccessKey,
		cfg.MinIO.BucketName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init minio client: %v", err)
	}

	// Init File repository
	fileRepo := repository.NewFileRepository(minioClient, cfg.MinIO.BucketName)

	// Init OpenAI service
	openAIService := ml.NewOpenAIService(cfg.OpenAI.BaseURL, cfg.OpenAI.ApiKey)

	// Init ML repository
	mlRepo := repository.NewMLRepository(rdb)

	// Init user repository and service
	userRepo := repository.NewUserRepository(queries)
	UserService := app.NewUserService(*userRepo)

	// Init user handler
	userHandler := handlers.NewUserHandler(UserService)

	// Init advertiser repository and service
	advertiserRepo := repository.NewAdvertiserRepository(queries)
	advertiserService := app.NewAdvertiserService(*advertiserRepo)

	// Init advertiser handler
	advertiserHandler := handlers.NewAdvertiserHandler(advertiserService)

	// Init time repository and service
	timeRepo := repository.NewTimeRepository(rdb)
	timeService := app.NewTimeService(*timeRepo)

	// Init time handler
	timeHandler := handlers.NewTimeHandler(timeService)

	// Init campaign repository and service
	campaignRepo := repository.NewCampaignRepository(queries, conn)
	campaignService := app.NewCampaignService(
		*campaignRepo,
		*advertiserRepo,
		*timeRepo,
		openAIService,
		*mlRepo,
		*fileRepo,
		cfg.MinIO.PublicHost)

	// Init campaign handler
	campaignHandler := handlers.NewCampaignHandler(campaignService)

	// Init ads repository and service
	adsRepo := repository.NewAdsRepository(queries)
	adsService := app.NewAdsService(*adsRepo, *userRepo, *campaignRepo)

	// Init ads handler
	adsHandler := handlers.NewAdsHandler(adsService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(jsonMiddleware)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Post("/clients/bulk", userHandler.CreateUsers)
	r.Get("/clients/{clientId}", userHandler.GetByID)

	r.Post("/advertisers/bulk", advertiserHandler.CreateAdvertisers)
	r.Get("/advertisers/{advertiserId}", advertiserHandler.GetByID)

	r.Post("/ml-scores", advertiserHandler.CreateUpdateMLScore)

	r.Post("/advertisers/{advertiserId}/campaigns", campaignHandler.CreateCampaign)
	r.Get("/advertisers/{advertiserId}/campaigns", campaignHandler.GetCampaignsByAdvertiserID)
	r.Get("/advertisers/{advertiserId}/campaigns/{campaignId}", campaignHandler.GetCampaignByID)
	r.Put("/advertisers/{advertiserId}/campaigns/{campaignId}", campaignHandler.UpdateCampaign)
	r.Delete("/advertisers/{advertiserId}/campaigns/{campaignId}", campaignHandler.DeleteCampaign)

	r.Post("/advertisers/{advertiserId}/campaigns/{campaignId}/picture", campaignHandler.SetCampaignPicture)

	r.Post("/advertisers/campaigns/generate", campaignHandler.GenerateAdText)

	r.Patch("/advertisers/campaigns/moderation", campaignHandler.SwitchModeration)

	r.Post("/ads/{adId}/click", adsHandler.Click)

	r.Post("/time/advance", timeHandler.SetCurrentDate)

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.ServerAddress,
			Handler: r,
		},
		db: conn,
	}, nil

}

func initMinIO(ctx context.Context, endpoint, accessKeyID, secretAccessKey, bucketName string) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			return nil, err
		}
		log.Println("Created MinIO bucket named", bucketName)

		// Allow anonymous read-only access
		policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)

		err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
		log.Println("Set public read policy for bucket", bucketName)
	} else {
		log.Println("Found existing MinIO bucket with name", bucketName)
	}

	return minioClient, nil
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func InitDB(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = 100 // Max 100 connections

	return pool, nil
}
