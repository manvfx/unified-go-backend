package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `json:"username" validate:"required,min=3,max=32"`
	Password       string             `json:"password" validate:"required,min=6"`
	Email          string             `bson:"email" validate:"required,email"`
	Verified       bool               `bson:"verified"`
	VerifiedAt     time.Time          `bson:"verified_at,omitempty"`
	LastLogin      time.Time          `bson:"last_login,omitempty"`
	LastLoginIP    string             `bson:"last_login_ip,omitempty"`
	LastLoginAgent string             `bson:"last_login_agent,omitempty"`
}
