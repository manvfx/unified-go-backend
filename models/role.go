package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `json:"name" validate:"required"`
	Permissions []string           `json:"permissions" validate:"required"`
}
