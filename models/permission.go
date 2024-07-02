package models

type Permission struct {
	Name string `json:"name" validate:"required"`
}
