package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Users *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Users: db.Collection("users"),
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	email = normalizeEmail(email)

	var user User
	err := r.Users.FindOne(ctx, bson.M{
		"emails.address": email,
	}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}

func (r *Repository) FindByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.Users.FindOne(ctx, bson.M{
		"_id": userID,
	}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, role, schoolID string) (*User, error) {
	now := time.Now()

	user := &User{
		ID: uuid.NewString(),
		Emails: []EmailRecord{
			{
				Address:  normalizeEmail(email),
				Verified: false,
			},
		},
		PasswordHash: passwordHash,
		Role:         role,
		Roles:        []string{role},
		SchoolID:     schoolID,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err := r.Users.InsertOne(ctx, user)
	return user, err
}