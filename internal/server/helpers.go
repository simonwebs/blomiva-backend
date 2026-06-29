package server

import (
	"log"
	"os"
	"strings"
)

func printBanner(logger *log.Logger) {
	logger.Println("===================================")
	logger.Println("      BLOMIVA BACKEND API")
	logger.Println("===================================")
}

func fallback(v string, def string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	return v
}

func validateConfig(cfg struct {
	JWTSecret string
	MongoURI  string
	DBName    string
}) {
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		panic("JWT_SECRET missing")
	}
	if strings.TrimSpace(cfg.MongoURI) == "" {
		panic("MONGO_URI missing")
	}
	if strings.TrimSpace(cfg.DBName) == "" {
		panic("DB_NAME missing")
	}
}