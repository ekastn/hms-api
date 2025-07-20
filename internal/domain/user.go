package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role defines the access level of a user.
type Role string

const (
	RoleAdmin        Role = "Admin"
	RoleDoctor       Role = "Doctor"
	RoleNurse        Role = "Nurse"
	RoleReceptionist Role = "Receptionist"
	RoleManagement   Role = "Management"
)

// @Description	User object
// @Description	Used for creating and updating users.
// @swagger:model
type UserEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"60d0fe4f53115a001f000001"`
	Name      string             `bson:"name" json:"name" example:"John Doe"`
	Email     string             `bson:"email,unique" json:"email" example:"john.doe@example.com"`
	Password  string             `bson:"password" json:"password" example:"securepassword123"` // Stores the hashed password
	Role      Role               `bson:"role" json:"role" example:"Admin"`
	IsActive  bool               `bson:"isActive" json:"isActive" example:true`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// UserDTO is a safe data transfer object for user info (without the password).
type UserDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	IsActive bool   `json:"isActive"`
}

// LoginRequest defines the structure for a login request.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse defines the structure for a successful login response.
type LoginResponse struct {
	Token string   `json:"token"`
	User  *UserDTO `json:"user"`
}

// @Description	Request body for creating a new user
// @swagger:model
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100" example:"Jane Doe"`
	Email    string `json:"email" validate:"required,email" example:"jane.doe@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"StrongPassword123"`
	Role     Role   `json:"role" validate:"required,oneof=Admin Doctor Nurse Receptionist Management" example:"Receptionist"`
}

// @Description	Request body for updating an existing user
// @swagger:model
type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" validate:"min=3,max=100" example:"Jane Doe"`
	Email    string `json:"email,omitempty" validate:"email" example:"jane.doe@example.com"`
	Role     Role   `json:"role,omitempty" validate:"oneof=Admin Doctor Nurse Receptionist Management" example:"Receptionist"`
	IsActive *bool  `json:"isActive,omitempty" example:true`
}

// @Description	Request body for changing user password
// @swagger:model
type ChangePasswordRequest struct {
	NewPassword     string `json:"newPassword" validate:"required,min=8" example:"VeryStrongNewPassword123"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword" example:"VeryStrongNewPassword123"`
}

func (u *UserEntity) ToDTO() *UserDTO {
    return &UserDTO{
        ID:       u.ID.Hex(),
        Name:     u.Name,
        Email:    u.Email,
        Role:     u.Role,
        IsActive: u.IsActive,
    }
}

func (u *UserDTO) ToEntity() *UserEntity {
	id, _ := primitive.ObjectIDFromHex(u.ID)
    return &UserEntity{
        ID:       id,
        Name:     u.Name,
        Email:    u.Email,
        Role:     u.Role,
        IsActive: u.IsActive,
    }
}

func (req *CreateUserRequest) ToEntity() *UserEntity {
	return &UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		IsActive: true, // New users are active by default
	}
}

func (req *UpdateUserRequest) ToEntity(existing *UserEntity) *UserEntity {
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Role != "" {
		existing.Role = req.Role
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	return existing
}
