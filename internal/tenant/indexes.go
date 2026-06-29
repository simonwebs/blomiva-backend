package tenant

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection(CollectionName)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "tenantId", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("tenantId_unique"),
		},
		{
			Keys: bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("slug_unique"),
		},
		{
			Keys: bson.D{{Key: "domain", Value: 1}},
			Options: options.Index().
				SetSparse(true).
				SetName("domain_sparse"),
		},
		{
			Keys: bson.D{{Key: "owner.userId", Value: 1}},
			Options: options.Index().
				SetName("owner_userId"),
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "active", Value: 1},
				{Key: "deleted", Value: 1},
			},
			Options: options.Index().
				SetName("status_active_deleted"),
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().
				SetName("createdAt_desc"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
