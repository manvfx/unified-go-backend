package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username"`
	Password       string             `bson:"password"`
	Email          string             `bson:"email"`
	Verified       bool               `bson:"verified"`
	VerifiedAt     time.Time          `bson:"verified_at,omitempty"`
	LastLogin      time.Time          `bson:"last_login,omitempty"`
	LastLoginIP    string             `bson:"last_login_ip,omitempty"`
	LastLoginAgent string             `bson:"last_login_agent,omitempty"`
}
