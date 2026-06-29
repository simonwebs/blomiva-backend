package contact

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Messages *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Messages: db.Collection("contact_messages"),
	}
}

func (r *Repository) Create(ctx context.Context, message *ContactMessage) error {
	_, err := r.Messages.InsertOne(ctx, message)
	return err
}

func (r *Repository) List(ctx context.Context, limit int64) ([]ContactMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	cursor, err := r.Messages.Find(
		ctx,
		bson.M{},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []ContactMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*ContactMessage, error) {
	var message ContactMessage

	err := r.Messages.FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return &message, err
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, status MessageStatus) error {
	_, err := r.Messages.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}},
	)

	return err
}
