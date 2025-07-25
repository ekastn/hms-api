package repository

import (
	"context"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MedicalRecordRepository struct {
	collection *mongo.Collection
}

func NewMedicalRecordRepository(coll *mongo.Collection) *MedicalRecordRepository {
	return &MedicalRecordRepository{
		collection: coll,
	}
}

func (r *MedicalRecordRepository) Create(ctx context.Context, record *domain.MedicalRecordEntity) (primitive.ObjectID, error) {
	now := time.Now()
	record.CreatedAt = now
	record.UpdatedAt = now

	res, err := r.collection.InsertOne(ctx, record)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *MedicalRecordRepository) FindAll(ctx context.Context) ([]*domain.MedicalRecordEntity, error) {
	return r.findRecords(ctx, bson.M{})
}

func (r *MedicalRecordRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.MedicalRecordEntity, error) {
	var record domain.MedicalRecordEntity
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &record, nil
}

func (r *MedicalRecordRepository) findRecords(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*domain.MedicalRecordEntity, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []*domain.MedicalRecordEntity
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *MedicalRecordRepository) GetByPatientID(ctx context.Context, patientID primitive.ObjectID) ([]*domain.MedicalRecordEntity, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}) // Sort by date in descending order
	return r.findRecords(ctx, bson.M{"patientId": patientID}, opts)
}

func (r *MedicalRecordRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.MedicalRecordEntity, error) {
	filter := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lte": end,
		},
		"isDeleted": bson.M{"$ne": true},
	}
	return r.findRecords(ctx, filter, options.Find().SetSort(bson.D{{Key: "date", Value: -1}}))
}

func (r *MedicalRecordRepository) Update(ctx context.Context, id primitive.ObjectID, record *domain.MedicalRecordEntity) error {
	record.UpdatedAt = time.Now()

	update := bson.M{
		"$set": record,
	}

	_, err := r.collection.UpdateByID(ctx, id, update)
	return err
}

func (r *MedicalRecordRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MedicalRecordRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"isDeleted": bson.M{"$ne": true}})
}
