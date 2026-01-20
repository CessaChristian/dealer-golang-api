package main

import (
	"dealer_golang_api/internal/config"
	"dealer_golang_api/internal/routes"
	"log"
	"os"

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// REGISTER ROUTES
	routes.RegisterRoutes(e, db)

	// RUN SERVER
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	log.Println("Server running on port " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
