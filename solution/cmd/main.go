package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/config"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/app"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/clients/bulk", userHandler.CreateUsers)
	r.Get("/clients/{clientId}", userHandler.GetByID)

	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
