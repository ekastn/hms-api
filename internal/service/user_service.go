package service

import (
	"context"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService handles business logic for user management.
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo}
}

// GetAllUsers retrieves all users.
func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.UserDTO, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var userDTOs []*domain.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, &domain.UserDTO{
			ID:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}
	return userDTOs, nil
}

// GetUserByID retrieves a single user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.UserDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, objID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Not found
	}

	return &domain.UserDTO{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

// UpdateUser updates a user's details (e.g., name, role).
func (s *UserService) UpdateUser(ctx context.Context, id string, user *domain.UserEntity) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.userRepo.Update(ctx, objID, user)
}

// DeactivateUser marks a user as inactive.
func (s *UserService) DeactivateUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.userRepo.Deactivate(ctx, objID)
}
