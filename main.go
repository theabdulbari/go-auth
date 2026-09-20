package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-auth/database"
	"go-auth/models"
	"go-auth/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system env vars")
	}

	// Connect DB + migrate
	database.Connect()

	// Start background cleanup goroutine with a stop channel
	stopCleanup := make(chan struct{})
	go cleanupExpiredTokens(stopCleanup)

	// Build router
	router := routes.SetupRouter()

	// Create http.Server with timeouts
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server in its own goroutine so main can wait for signals
	go func() {
		log.Printf("Server listening on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for SIGINT (Ctrl+C) or SIGTERM (docker stop / k8s)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	// Give in-flight requests 10 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	} else {
		log.Println("HTTP server stopped cleanly")
	}

	// Stop the background goroutine
	close(stopCleanup)
	log.Println("Bye")
}

// cleanupExpiredTokens runs on a 6-hour ticker until stopCleanup is closed.
func cleanupExpiredTokens(stop <-chan struct{}) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	// Run once at startup so we don't wait 6 hours for the first sweep
	pruneExpiredTokens()

	for {
		select {
		case <-ticker.C:
			pruneExpiredTokens()
		case <-stop:
			log.Println("Cleanup goroutine stopped")
			return
		}
	}
}

func pruneExpiredTokens() {
	res := database.DB.
		Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{})

	if res.Error != nil {
		log.Printf("token cleanup failed: %v", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("Pruned %d expired refresh tokens", res.RowsAffected)
	}
}