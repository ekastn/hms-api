package repository

import (
	"context"

	"github.com/ekastn/hms-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DoctorRepository struct {
	coll *mongo.Collection
}

func NewDoctorRepository(coll *mongo.Collection) *DoctorRepository {
	return &DoctorRepository{coll}
}

func (r *DoctorRepository) GetByEmail(ctx context.Context, email string) (*domain.DoctorEntity, error) {
	var doctor domain.DoctorEntity
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&doctor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &doctor, nil
}

func (r *DoctorRepository) GetByPhone(ctx context.Context, phone string) (*domain.DoctorEntity, error) {
	var doctor domain.DoctorEntity
	err := r.coll.FindOne(ctx, bson.M{"phone": phone}).Decode(&doctor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &doctor, nil
}

func (r *DoctorRepository) GetByName(ctx context.Context, name string) (*domain.DoctorEntity, error) {
	var doctor domain.DoctorEntity
	err := r.coll.FindOne(ctx, bson.M{"name": name}).Decode(&doctor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &doctor, nil
}

func (r *DoctorRepository) Create(ctx context.Context, doctor *domain.DoctorEntity) (primitive.ObjectID, error) {
	res, err := r.coll.InsertOne(ctx, doctor)
	return res.InsertedID.(primitive.ObjectID), err
}

func (r *DoctorRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.DoctorEntity, error) {
	var doctor domain.DoctorEntity
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doctor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &doctor, nil
}

func (r *DoctorRepository) GetAll(ctx context.Context) ([]*domain.DoctorEntity, error) {
	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var doctors []*domain.DoctorEntity
	for cur.Next(ctx) {
		var doctor domain.DoctorEntity
		if err := cur.Decode(&doctor); err != nil {
			return nil, err
		}
		doctors = append(doctors, &doctor)
	}
	return doctors, nil
}

func (r *DoctorRepository) Update(ctx context.Context, id primitive.ObjectID, doctor *domain.DoctorEntity) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": doctor})
	return err
}

func (r *DoctorRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *DoctorRepository) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}
