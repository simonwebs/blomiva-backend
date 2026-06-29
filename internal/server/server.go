package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blomiva-backend/internal/config"
	"blomiva-backend/internal/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	cfg    config.Config
	router *gin.Engine
	db     *mongo.Database
	logger *log.Logger
	server *http.Server
}

func New(cfg config.Config, logger *log.Logger) (*App, error) {
	if logger == nil {
		logger = log.New(os.Stdout, "[blomiva] ", log.LstdFlags|log.LUTC)
	}

	printBanner(logger)
	validateConfig(cfg)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	db, err := database.ConnectMongo(ctx, cfg.MongoURI, cfg.DBName)
	if err != nil {
		return nil, fmt.Errorf("mongo connect failed: %w", err)
	}

	router := gin.New()

	// middleware (from middleware.go)
	router.Use(RecoveryMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(RequestLogger(logger))
	router.Use(CORSMiddleware(cfg))

	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("trusted proxies error: %w", err)
	}

	// routes (from routes.go)
	RegisterRoutes(router, db, cfg)

	port := fallback(cfg.Port, "8082")

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	return &App{
		cfg:    cfg,
		router: router,
		db:     db,
		logger: logger,
		server: srv,
	}, nil
}

func (a *App) Run() error {
	a.logger.Println("server: starting...")

	errChan := make(chan error, 1)

	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		errChan <- nil
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		a.logger.Println("shutdown signal:", sig)
	case err := <-errChan:
		if err != nil {
			return err
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.server.Shutdown(ctx)
}