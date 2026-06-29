package auth

import "time"

type EmailRecord struct {
	Address  string `bson:"address" json:"address"`
	Verified bool   `bson:"verified" json:"verified"`
}

type User struct {
	ID           string        `bson:"_id" json:"id"`
	Emails       []EmailRecord `bson:"emails" json:"emails"`
	PasswordHash string        `bson:"passwordHash,omitempty" json:"-"`
	Role         string        `bson:"role" json:"role"`
	Roles        []string      `bson:"roles" json:"roles"`
	SchoolID     string        `bson:"schoolId,omitempty" json:"schoolId,omitempty"`
	IsActive     bool          `bson:"isActive" json:"isActive"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
	SchoolID string `json:"schoolId"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthUserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	SchoolID string `json:"schoolId,omitempty"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  AuthUserResponse `json:"user"`
}