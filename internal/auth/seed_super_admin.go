package auth

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedSuperAdmin(repo *Repository) {
	email := strings.TrimSpace(os.Getenv("SUPER_ADMIN_EMAIL"))
	password := strings.TrimSpace(os.Getenv("SUPER_ADMIN_PASSWORD"))

	if email == "" || password == "" {
		log.Println("super-admin: missing env config, skipping seed")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existing, err := repo.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		log.Println("super-admin seed error:", err)
		return
	}

	if existing != nil {
		log.Println("super-admin already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("super-admin hash error:", err)
		return
	}

	_, err = repo.CreateUser(
		ctx,
		normalizeEmail(email),
		string(hash),
		"super-admin",
		"",
	)

	if err != nil {
		log.Println("super-admin create error:", err)
		return
	}

	log.Println("super-admin created successfully")
}
