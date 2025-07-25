package service

import (
	"context"
	"errors"
	"time"

	"github.com/ekastn/hms-api/internal/domain"
	"github.com/ekastn/hms-api/internal/repository"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles user authentication and JWT generation.
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Login verifies a user's credentials and returns a JWT and user DTO.
func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user == nil || !user.IsActive {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	userDTO := &domain.UserDTO{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return &domain.LoginResponse{
		Token: token,
		User:  userDTO,
	}, nil
}

// generateJWT creates a new JWT for a given user.
func (s *AuthService) generateJWT(user *domain.UserEntity) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.Hex(),
		"role":  user.Role,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
