package core

import (
	"booking-app/database"
	"booking-app/model"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


var eventCollection *mongo.Collection = database.OpenCollection(database.Client, "event")


func PublishEvent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		eventCtx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		var event model.PublishEvent

		if err := ctx.BindJSON(&event); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
