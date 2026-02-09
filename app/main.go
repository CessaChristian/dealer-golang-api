package main

import (
	"context"
	"dealer_golang_api/internal/config"
	"dealer_golang_api/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// LOAD ENV
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ .env file not found, using system environment variables")
	}

	// INIT DATABASE (PGXPOOL)
	db := config.InitDB()
	defer db.Close()

	// INIT ECHO
	e := echo.New()

	// MIDDLEWARE
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} | ${latency_human} | ${remote_ip} | ${method} ${uri} | ${error}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// REGISTER ROUTES
	routes.RegisterRoutes(e, db)

	// RUN SERVER
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	// Start server in goroutine
	go func() {
		log.Println("Server running on port " + port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: ", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C / SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give 10 seconds for in-flight requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server stopped gracefully")
}
