package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/app"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/config"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/handlers"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

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

	// Init campaign repository and service
	campaignRepo := repository.NewCampaignRepository(queries, conn)
	campaignService := app.NewCampaignService(*campaignRepo)

	// Init campaign handler
	campaignHandler := handlers.NewCampaignHandler(campaignService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

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

	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
