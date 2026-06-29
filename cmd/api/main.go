package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"blomiva-backend/internal/config"
	"blomiva-backend/internal/database"
	"blomiva-backend/internal/server"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	logger := log.New(os.Stdout, "[blomiva] ", log.LstdFlags|log.LUTC)

	// Load .env (optional)
	if err := godotenv.Load(); err != nil {
		logger.Println("config: .env not found, using system environment variables")
	}

	cfg := config.Load()

	// ================= DB CONNECTION + HEALTH CHECK =================
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Println("database: connecting...")

	db, err := database.ConnectMongo(ctx, cfg.MongoURI, cfg.DBName)
	if err != nil {
		logger.Fatalf("database: connection failed: %v", err)
	}

	if err := db.Client().Ping(ctx, readpref.Primary()); err != nil {
		logger.Fatalf("database: ping failed: %v", err)
	}

	logger.Println("database: connected successfully")

	// ================= SERVER START =================
	app, err := server.New(cfg, logger, db)
	if err != nil {
		logger.Fatalf("server: init failed: %v", err)
	}

	if err := app.Run(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Println("server: gracefully stopped")
			return
		}
		logger.Fatalf("server: failed: %v", err)
	}
}
