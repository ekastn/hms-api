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

//	@Description	User object
//	@Description	Used for creating and updating users.
//	@swagger:model
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
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
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
