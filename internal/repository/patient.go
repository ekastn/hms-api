package repository

import (
	"context"

	"github.com/ekastn/hms-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PatientRepository struct {
	coll *mongo.Collection
}

func NewPatientRepository(coll *mongo.Collection) *PatientRepository {
	return &PatientRepository{coll}
}

func (r *PatientRepository) Create(ctx context.Context, patient *domain.PatientEntity) (primitive.ObjectID, error) {
	res, err := r.coll.InsertOne(ctx, patient)
	return res.InsertedID.(primitive.ObjectID), err
}

func (r *PatientRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.PatientEntity, error) {
	var patient domain.PatientEntity
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&patient)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) GetByEmail(ctx context.Context, email string) (*domain.PatientEntity, error) {
	var patient domain.PatientEntity
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&patient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) GetByPhone(ctx context.Context, phone string) (*domain.PatientEntity, error) {
	var patient domain.PatientEntity
	err := r.coll.FindOne(ctx, bson.M{"phone": phone}).Decode(&patient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) GetByName(ctx context.Context, name string) (*domain.PatientEntity, error) {
	var patient domain.PatientEntity
	err := r.coll.FindOne(ctx, bson.M{"name": name}).Decode(&patient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{"isDeleted": bson.M{"$ne": true}})
}

func (r *PatientRepository) GetAll(ctx context.Context) ([]*domain.PatientEntity, error) {
	cur, err := r.coll.Find(ctx, bson.M{"isDeleted": bson.M{"$ne": true}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var patients []*domain.PatientEntity
	for cur.Next(ctx) {
		var patient domain.PatientEntity
		if err := cur.Decode(&patient); err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}
	return patients, nil
}

func (r *PatientRepository) Update(ctx context.Context, id primitive.ObjectID, patient *domain.PatientEntity) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": patient,
	})

	return err
}

func (r *PatientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
