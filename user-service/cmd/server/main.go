package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lonaas1337-sudo/vinylvault/user-service/internal/config"
	"github.com/lonaas1337-sudo/vinylvault/user-service/internal/handler"
	"github.com/lonaas1337-sudo/vinylvault/user-service/internal/repository"
)

func main() {
	cfg := config.Load()

	repo, err := repository.NewUserRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to conntect to database: %v", err)
	}

	defer repo.Close()
	handler.SetRepository(repo)
	fmt.Println("Successfully connected to PostgreSQL!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", handler.HealthHandler)
	r.Post("/users/register", handler.RegisterHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in the background
	go func() {
		fmt.Printf("User service starting on http://localhost:%s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Make a channel to listen to shutdown signal (Ctrl+C or kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the signal
	<-quit
	log.Println("Shutting down user service...")

	// Give 10 seconds to finish current requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shuwdown server with context (you have 10 secods to finish)
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("User service stopped gracefully")
}
