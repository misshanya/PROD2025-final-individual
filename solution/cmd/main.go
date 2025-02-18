package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "gitlab.prodcontest.ru/2025-final-projects-back/misshanya/docs"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/app"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/config"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/handlers"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/ml"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

// @title			PROD Backend 2025 Advertising Platform API
// @version		1.0
// @description	API для управления данными клиентов, рекламодателей, рекламными кампаниями, показом объявлений, статистикой и управлением "текущим днём" в системе.
// @license.name	GPL 3.0
func main() {
	cfg := config.NewConfig()

	ctx := context.Background()

	// Init db connection
	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(ctx)

	// Init SQL queries
	queries := storage.New(conn)

	// Init redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to ping redis: %v", err)
	}

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
	campaignService := app.NewCampaignService(*campaignRepo, *timeRepo, openAIService, *mlRepo)

	// Init campaign handler
	campaignHandler := handlers.NewCampaignHandler(campaignService)

	// Init ads repository and service
	adsRepo := repository.NewAdsRepository(queries)
	adsService := app.NewAdsService(*adsRepo, *userRepo, *campaignRepo)

	// Init ads handler
	adsHandler := handlers.NewAdsHandler(adsService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
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

	r.Post("/advertisers/campaigns/generate", campaignHandler.GenerateAdText)

	r.Patch("/advertisers/campaigns/moderation", campaignHandler.SwitchModeration)

	r.Post("/ads/{adId}/click", adsHandler.Click)

	r.Post("/time/advance", timeHandler.SetCurrentDate)

	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
