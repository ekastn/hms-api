package service

import (
	"context"
	"errors"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"github.com/ekastn/hms-api/internal/utils"
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
		userDTOs = append(userDTOs, user.ToDTO())
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

    return user.ToDTO(), nil
}

// UpdateUser updates a user's details (e.g., name, role).
func (s *UserService) UpdateUser(ctx context.Context, id string, req *domain.UpdateUserRequest) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	existingUser, err := s.userRepo.GetByID(ctx, objID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	updatedUser := req.ToEntity(existingUser)

	return s.userRepo.Update(ctx, objID, updatedUser)
}

// ChangeUserPassword updates a user's password.
func (s *UserService) ChangeUserPassword(ctx context.Context, id string, newPassword string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	existingUser, err := s.userRepo.GetByID(ctx, objID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	existingUser.Password = string(hashedPassword)

	return s.userRepo.Update(ctx, objID, existingUser)
}

// DeactivateUser marks a user as inactive.
func (s *UserService) DeactivateUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.userRepo.Deactivate(ctx, objID)
}

// CreateUser creates a new user with a hashed password.
func (s *UserService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (string, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user := req.ToEntity()
	user.Password = string(hashedPassword)

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}
