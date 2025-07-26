package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener-go-backend/internal/handler"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/router"
	"url-shortener-go-backend/internal/service"

	"github.com/joho/godotenv"
)

func main() {
err1 := godotenv.Load(".env")

	if err1 != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	apiKey := os.Getenv("SERVICE_ROLE")
	
	apiURL := os.Getenv("DB_API_URL")+"/rest/v1"

	if apiKey == "" || apiURL == "" {
		log.Fatal("SERVICE_ROLE or DB_API_URL environment variable is missing")
	}

	supabase, err := repository.NewSupabaseRepository(apiURL, apiKey)
	if err != nil {
		log.Fatalf("failed to create supabase repository: %v", err)
	}

	urlRepo := repository.NewURLRepository(supabase)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)
	server := router.NewAPIServer(":8080", urlHandler) //change in prod

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Received shutdown signal...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}
