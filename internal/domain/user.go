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

// UserEntity represents a staff member with login credentials.
type UserEntity struct {
ID        primitive.ObjectID `bson:"_id,omitempty"`
Name      string             `bson:"name"`
Email     string             `bson:"email,unique"`
Password  string             `bson:"password"` // Stores the hashed password
Role      Role               `bson:"role"`
IsActive  bool               `bson:"isActive"`
CreatedAt time.Time          `bson:"createdAt"`
UpdatedAt time.Time          `bson:"updatedAt"`
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
