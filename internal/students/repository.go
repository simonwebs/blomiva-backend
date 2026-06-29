package students

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrStudentNotFound = errors.New("student not found")

type Repository interface {
	Create(ctx context.Context, student *Student) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("students"),
	}
}

func (r *MongoRepository) Create(ctx context.Context, student *Student) error {
	_, err := r.collection.InsertOne(ctx, student)
	return err
}

func NewObjectIDFromHex(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
