package repository

import (
	"context"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentRepository struct {
	coll *mongo.Collection
}

func NewAppointmentRepository(coll *mongo.Collection) *AppointmentRepository {
	return &AppointmentRepository{coll}
}

func (r *AppointmentRepository) Create(ctx context.Context, appointment *domain.AppointmentEntity) (primitive.ObjectID, error) {
	now := time.Now()
	appointment.CreatedAt = now
	appointment.UpdatedAt = now

	res, err := r.coll.InsertOne(ctx, appointment)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.AppointmentEntity, error) {
	var appointment domain.AppointmentEntity
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&appointment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) GetAll(ctx context.Context) ([]*domain.AppointmentEntity, error) {
	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var appointments []*domain.AppointmentEntity
	for cur.Next(ctx) {
		var appointment domain.AppointmentEntity
		if err := cur.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *AppointmentRepository) Update(ctx context.Context, id primitive.ObjectID, appointment *domain.AppointmentEntity) error {
	appointment.UpdatedAt = time.Now()

	_, err := r.coll.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": appointment},
	)

	return err
}

func (r *AppointmentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *AppointmentRepository) GetByDoctorAndDateRange(
	ctx context.Context,
	doctorID primitive.ObjectID,
	start, end time.Time,
) ([]*domain.AppointmentEntity, error) {
	filter := bson.M{
		"doctorId": doctorID,
		"dateTime": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)


	var appointments []*domain.AppointmentEntity
	for cur.Next(ctx) {
		var appointment domain.AppointmentEntity
		if err := cur.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}
