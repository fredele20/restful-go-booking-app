package model

type PublishEvent struct {
	Name           string `json:"name" validate:"required,min=2"`
	Venue          string `json:"venue" validate:"required,min=3"`
	TicketQuantity int    `json:"quantity" validate:"required,min=10"`
}
