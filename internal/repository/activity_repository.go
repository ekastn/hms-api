package repository

import (
	"context"

	"github.com/ekastn/hms-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActivityRepository struct {
	coll *mongo.Collection
}

func NewActivityRepository(coll *mongo.Collection) *ActivityRepository {
	return &ActivityRepository{coll}
}

func (r *ActivityRepository) Create(ctx context.Context, activity *domain.ActivityEntity) error {
	_, err := r.coll.InsertOne(ctx, activity)
	return err
}

func (r *ActivityRepository) GetRecent(ctx context.Context, limit int) ([]*domain.ActivityEntity, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))
	cur, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var activities []*domain.ActivityEntity
	if err := cur.All(ctx, &activities); err != nil {
		return nil, err
	}

	return activities, nil
}
