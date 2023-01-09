package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PublishEvent struct {
	ID             primitive.ObjectID `json:"_id"`
	Name           string             `json:"name" validate:"required,min=2"`
	Venue          string             `json:"venue" validate:"required,min=3"`
	TicketQuantity int                `json:"quantity" validate:"required,min=10"`
}

type EventBooking struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	NumOfTickets int                `json:"numOfTickets" validate:"required,min=1"`
}
