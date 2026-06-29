package profiles

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Profiles  *mongo.Collection
	Users     *mongo.Collection
	AuditLogs *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *Repository {
	return &Repository{
		Profiles:  db.Collection("profiles"),
		Users:     db.Collection("users"),
		AuditLogs: db.Collection("audit_logs"),
	}
}

func NewRepository(db *mongo.Database) *Repository {
	return NewMongoRepository(db)
}

func (r *Repository) FindProfileByOwnerID(ctx context.Context, ownerID string) (*Profile, error) {
	var profile Profile

	err := r.Profiles.FindOne(ctx, bson.M{
		"ownerId": ownerID,
	}).Decode(&profile)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &profile, err
}

func (r *Repository) FindProfileBySlug(ctx context.Context, slug string) (*Profile, error) {
	var profile Profile

	err := r.Profiles.FindOne(ctx, bson.M{
		"slug":      slug,
		"status":    "active",
		"isBlocked": bson.M{"$ne": true},
	}).Decode(&profile)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &profile, err
}

func (r *Repository) FindUserByID(ctx context.Context, userID string) (*User, error) {
	var user User

	err := r.Users.FindOne(ctx, bson.M{
		"_id": userID,
	}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.Users.FindOne(ctx, bson.M{
		"emails.address": email,
	}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &user, err
}

func (r *Repository) InsertProfile(ctx context.Context, profile *Profile) error {
	_, err := r.Profiles.InsertOne(ctx, profile)
	return err
}

func (r *Repository) UpdateProfileByOwnerID(ctx context.Context, ownerID string, set bson.M) error {
	if set == nil {
		set = bson.M{}
	}

	set["updatedAt"] = time.Now()

	_, err := r.Profiles.UpdateOne(
		ctx,
		bson.M{"ownerId": ownerID},
		bson.M{"$set": set},
	)

	return err
}

func (r *Repository) UnsetProfileCustomKey(ctx context.Context, ownerID string, key string) error {
	_, err := r.Profiles.UpdateOne(
		ctx,
		bson.M{"ownerId": ownerID},
		bson.M{
			"$unset": bson.M{
				"custom." + key: "",
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
			},
		},
	)

	return err
}

func (r *Repository) UpdateUser(ctx context.Context, userID string, set bson.M) error {
	if set == nil {
		set = bson.M{}
	}

	set["updatedAt"] = time.Now()

	_, err := r.Users.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": set},
	)

	return err
}

func (r *Repository) DeleteUserAndProfile(ctx context.Context, ownerID string) error {
	_, err := r.Profiles.DeleteOne(ctx, bson.M{
		"ownerId": ownerID,
	})

	if err != nil {
		return err
	}

	_, err = r.Users.DeleteOne(ctx, bson.M{
		"_id": ownerID,
	})

	return err
}

func (r *Repository) CreateAuditLog(ctx context.Context, action string, data bson.M) error {
	if data == nil {
		data = bson.M{}
	}

	now := time.Now()

	data["action"] = action
	data["timestamp"] = now
	data["createdAt"] = now

	_, err := r.AuditLogs.InsertOne(ctx, data)
	return err
}