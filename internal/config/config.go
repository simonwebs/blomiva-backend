package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env             string
	Port            string
	MongoURI        string
	DBName          string
	JWTSecret       string
	JWTExpiresHours int

	AppURL string
}

// Load builds application config safely
func Load() Config {
	env := getEnv("APP_ENV", "development")

	mongoURI := firstEnv("MONGO_URI", "MONGO_URL")

	if mongoURI == "" && env == "production" {
		log.Fatal("missing required env: MONGO_URI")
	}

	cfg := Config{
		Env:             env,
		Port:            getEnv("PORT", "8080"),
		MongoURI:        defaultIfEmpty(mongoURI, "mongodb://localhost:27017"),
		DBName:          getEnv("DB_NAME", defaultDB(env)),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		JWTExpiresHours: getEnvAsInt("JWT_EXPIRES_HOURS", 24),
		AppURL:          getEnv("APP_URL", "http://localhost:8080"),
	}

	validate(cfg)

	return cfg
}

// ================= helpers =================

func getEnv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func getEnvAsInt(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}

	return n
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}

func defaultIfEmpty(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}

func defaultDB(env string) string {
	switch env {
	case "production":
		return "blomiva"
	case "staging":
		return "blomiva_staging"
	default:
		return "blomiva_dev"
	}
}

// ================= validation =================

func validate(cfg Config) {
	if cfg.Env == "production" && cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required in production")
	}

	if cfg.JWTExpiresHours > 72 {
		log.Println("warning: JWT_EXPIRES_HOURS is very high")
	}
}
