package tenant

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrTenantNotFound = errors.New("tenant not found")

type Repository interface {
	Create(ctx context.Context, t *Tenant) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*Tenant, error)
	FindByTenantID(ctx context.Context, tenantID string) (*Tenant, error)
	FindBySlug(ctx context.Context, slug string) (*Tenant, error)
	List(ctx context.Context, page, limit int64) ([]Tenant, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*Tenant, error)
	SoftDelete(ctx context.Context, id primitive.ObjectID) error
}

type MongoRepository struct {
	col *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{col: db.Collection(CollectionName)}
}

func (r *MongoRepository) Create(ctx context.Context, t *Tenant) error {
	_, err := r.col.InsertOne(ctx, t)
	return err
}

func (r *MongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Tenant, error) {
	var t Tenant

	err := r.col.FindOne(ctx, bson.M{
		"_id":     id,
		"deleted": false,
	}).Decode(&t)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrTenantNotFound
	}

	return &t, err
}

func (r *MongoRepository) FindByTenantID(ctx context.Context, tenantID string) (*Tenant, error) {
	var t Tenant

	err := r.col.FindOne(ctx, bson.M{
		"tenantId": tenantID,
		"deleted":  false,
	}).Decode(&t)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrTenantNotFound
	}

	return &t, err
}

func (r *MongoRepository) FindBySlug(ctx context.Context, slug string) (*Tenant, error) {
	var t Tenant

	err := r.col.FindOne(ctx, bson.M{
		"slug":    slug,
		"deleted": false,
	}).Decode(&t)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrTenantNotFound
	}

	return &t, err
}

func (r *MongoRepository) List(ctx context.Context, page, limit int64) ([]Tenant, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := bson.M{"deleted": false}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var items []Tenant
	if err := cur.All(ctx, &items); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *MongoRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*Tenant, error) {
	update["updatedAt"] = time.Now().UTC()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var t Tenant
	err := r.col.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id":     id,
			"deleted": false,
		},
		bson.M{"$set": update},
		opts,
	).Decode(&t)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrTenantNotFound
	}

	return &t, err
}

func (r *MongoRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now().UTC()

	res, err := r.col.UpdateOne(ctx,
		bson.M{
			"_id":     id,
			"deleted": false,
		},
		bson.M{"$set": bson.M{
			"status":    StatusArchived,
			"active":    false,
			"deleted":   true,
			"archived":  true,
			"deletedAt": now,
			"updatedAt": now,
		}},
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrTenantNotFound
	}

	return nil
}
