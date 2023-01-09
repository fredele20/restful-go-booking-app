package core

import (
	"booking-app/database"
	"booking-app/model"
	"booking-app/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var eventCollection *mongo.Collection = database.OpenCollection(database.Client, "event")
var bookedTicketCollection *mongo.Collection = database.OpenCollection(database.Client, "bookedTickets")

func PublishEvent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		eventCtx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var event model.PublishEvent

		if err := ctx.BindJSON(&event); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if validationErr := utils.ValidateInput(event, ctx); validationErr != nil {
			return
		}

		_, insertErr := eventCollection.InsertOne(eventCtx, event)
		if insertErr != nil {
			msg := fmt.Sprintf("unsuccessful, event item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, event)
	}
}

// I have not tested this yet, and I am expecting a log of errors
func BookEventTicket() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bookingCtx, _ := context.WithTimeout(context.Background(), 100 * time.Second)
		eventCtx, _ := context.WithTimeout(context.Background(), 100 * time.Second)
		var booking model.EventBooking
		var publishedEvent model.PublishEvent

		if err := ctx.BindJSON(&bookingCtx); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if validationErr := utils.ValidateInput(booking, ctx); validationErr != nil {
			return
		}

		count, err := userCollection.CountDocuments(eventCtx, bson.M{"_id": publishedEvent.ID})
		// defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the id"})
			return
		}

		if count <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no event with the given ID found"})
			return
		}
		
		upsert := true
		filter := bson.M{"_id": publishedEvent.ID}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		eventCollection.UpdateByID(
			eventCtx, 
			filter,
			bson.D{
				{"$set", bson.E{"quantity", publishedEvent.TicketQuantity - booking.NumOfTickets},
				},
			},
			&opt,
		)

		_, bookingErr := bookedTicketCollection.InsertOne(bookingCtx, booking)
		if bookingErr != nil {
			msg := fmt.Sprintf("unsuccessful, unable to perform booking operation")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// defer cancel()
		ctx.JSON(http.StatusOK, booking)
	}
}
