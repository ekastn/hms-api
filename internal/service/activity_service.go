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

// GetAllActivities retrieves all activities.
func (s *ActivityService) GetAllActivities(ctx context.Context) ([]*domain.Activity, error) {
	activities, err := s.activityRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var activityDTOs []*domain.Activity
	for _, activity := range activities {
		activityDTOs = append(activityDTOs, &domain.Activity{
			ID:          activity.ID.Hex(),
			Type:        string(activity.Type),
			Title:       activity.Title,
			Description: activity.Description,
			Timestamp:   activity.Timestamp,
		})
	}
	return activityDTOs, nil
}
