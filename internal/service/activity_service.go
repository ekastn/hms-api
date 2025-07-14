package service

import (
	"context"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
)

type ActivityService struct {
	activityRepo *repository.ActivityRepository
}

func NewActivityService(activityRepo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{activityRepo}
}

func (s *ActivityService) CreateActivity(ctx context.Context, activityType domain.ActivityType, title, description string) error {
	activity := &domain.ActivityEntity{
		Type:        activityType,
		Title:       title,
		Description: description,
		Timestamp:   time.Now(),
	}

	return s.activityRepo.Create(ctx, activity)
}
