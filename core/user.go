package core

import (
	"booking-app/database"
	"booking-app/model"
	"booking-app/utils"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

var validate = validator.New()


func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userCtx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		var user model.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(userCtx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the email"})
			return
		}

		count, err = userCollection.CountDocuments(userCtx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the phone number"})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "a user with the email or phone number already exist"})
			return
		}

		password := utils.HashPassword(*user.Password)
		user.Password = &password

		user.ID = primitive.NewObjectID()
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		
	}
}