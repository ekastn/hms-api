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

func (r *AppointmentRepository) GetUpcomingAppointments(ctx context.Context, limit int) ([]*domain.UpcomingAppointment, error) {
	now := time.Now()
	end := now.Add(7 * 24 * time.Hour) // Next 7 days

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"dateTime": bson.M{
					"$gte": now,
					"$lte": end,
				},
				"status": bson.M{"$in": []string{"Scheduled", "Confirmed"}},
			},
		},
		{"$sort": bson.M{"date": 1}},
		{"$limit": limit},
		{
			"$lookup": bson.M{
				"from":         "patients",
				"localField":   "patientId",
				"foreignField": "_id",
				"as":           "patient",
			},
		},
		{"$unwind": "$patient"},
		{
			"$lookup": bson.M{
				"from":         "doctors",
				"localField":   "doctorId",
				"foreignField": "_id",
				"as":           "doctor",
			},
		},
		{"$unwind": "$doctor"},
		{
			"$project": bson.M{
				"_id":         1,
				"patientName": "$patient.name",
				"doctorName":  "$doctor.name",
				"dateTime":    1,
				"status":      1,
			},
		},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var appointments []*domain.UpcomingAppointment
	if err := cursor.All(ctx, &appointments); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *AppointmentRepository) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}

// GetAppointmentsCount is kept for backward compatibility
func (r *AppointmentRepository) GetAppointmentsCount(ctx context.Context) (int64, error) {
	return r.Count(ctx)
}

func (r *AppointmentRepository) GetByDoctorAndDateRange(
	ctx context.Context,
	doctorID primitive.ObjectID,
	start, end time.Time,
) ([]*domain.AppointmentEntity, error) {
	filter := bson.M{
		"doctorId": doctorID,
		"startTime": bson.M{
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
	if err := cur.All(ctx, &appointments); err != nil {
		return nil, err
	}

	return appointments, nil
}

// GetByPatientID retrieves appointments for a specific patient
func (r *AppointmentRepository) GetByPatientID(ctx context.Context, patientID primitive.ObjectID) ([]*domain.AppointmentEntity, error) {
	filter := bson.M{"patientId": patientID}
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}})

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var appointments []*domain.AppointmentEntity
	if err := cur.All(ctx, &appointments); err != nil {
		return nil, err
	}

	return appointments, nil
}

// GetRecentPatientsByDoctorID returns a list of recent patients for a doctor
func (r *AppointmentRepository) GetRecentPatientsByDoctorID(ctx context.Context, doctorID primitive.ObjectID, limit int) ([]primitive.ObjectID, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"doctorId": doctorID,
			},
		},
		{
			"$sort": bson.M{"startTime": -1},
		},
		{
			"$group": bson.M{
				"_id":       "$patientId",
				"lastVisit": bson.M{"$first": "$startTime"},
			},
		},
		{
			"$sort": bson.M{"lastVisit": -1},
		},
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	patientIDs := make([]primitive.ObjectID, len(results))
	for i, result := range results {
		patientIDs[i] = result.ID
	}

	return patientIDs, nil
}
