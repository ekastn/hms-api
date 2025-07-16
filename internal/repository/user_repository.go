
package repository

import (
	"context"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(coll *mongo.Collection) *UserRepository {
	return &UserRepository{coll}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.UserEntity) (primitive.ObjectID, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true

	res, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.UserEntity, error) {
	var user domain.UserEntity
	err := r.coll.FindOne(ctx, bson.M{"email": email, "isActive": true}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.UserEntity, error) {
	var user domain.UserEntity
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "isActive": true}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.UserEntity, error) {
	cur, err := r.coll.Find(ctx, bson.M{"isActive": true})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []*domain.UserEntity
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id primitive.ObjectID, user *domain.UserEntity) error {
	user.UpdatedAt = time.Now()
	update := bson.M{
		"$set": user,
	}
	_, err := r.coll.UpdateByID(ctx, id, update)
	return err
}

// Deactivate sets a user to inactive instead of deleting them.
func (r *UserRepository) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"isActive":  false,
			"updatedAt": time.Now(),
		},
	}
	_, err := r.coll.UpdateByID(ctx, id, update)
	return err
}
