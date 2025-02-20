package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gitlab.prodcontest.ru/2025-final-projects-back/misshanya/docs"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/config"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/server"
)

// @title			PROD Backend 2025 Advertising Platform API
// @version		1.0
// @description	API для управления данными клиентов, рекламодателей, рекламными кампаниями, показом объявлений, статистикой и управлением "текущим днём" в системе.
// @license.name	GPL 3.0
func main() {
	cfg := config.NewConfig()

	ctx := context.Background()

	server, err := server.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		if err := server.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down, bye bye...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.HttpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown: %v", err)
	}

	server.DB.Close()
}
