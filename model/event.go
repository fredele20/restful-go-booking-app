package model

type PublishEvent struct {
	Name     string `json:"name" validate:"required"`
	Venue    string `json:"venue" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}
