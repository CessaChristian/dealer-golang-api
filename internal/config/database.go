package config

import (
	"context"
	"dealer_golang_api/utils"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	dsn, err := utils.GetDecodedDBURL()
	if err != nil {
		log.Fatal("Failed to load DB URL:", err)
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Failed to parse DB URL:", err)
	}

	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}

	// Test koneksi
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to Database (using Base64 secure env)")
	return db
}
