package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessGroup struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `json:"name" validate:"required"`
	Roles       []string           `json:"roles"`
	Permissions []string           `json:"permissions"`
}
