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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

// func Count(ctx context.Context, input string) (count int64, err error) {
// 	count, err = userCollection.CountDocuments(ctx, bson.M{"email": input})
// 	if err != nil {
// 		log.Panic(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the email"})
// 		return count, err
// 	}
// 	return count, nil
// }

func checkExistingUser(ctx context.Context, field, value string) (int64, error) {
	count, err := userCollection.CountDocuments(ctx, bson.M{field: value})
	if err != nil {
		log.Panic(err)
		fmt.Printf("error checking for %v", field)
		return count, err
	}

	return count, nil
}

// count, err := userCollection.CountDocuments(userCtx, bson.M{"email": user.Email})
// defer cancel()
// if err != nil {
// 	log.Panic(err)
// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the email"})
// 	return
// }

// count, err = userCollection.CountDocuments(userCtx, bson.M{"phone": user.Phone})
// defer cancel()
// if err != nil {
// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the phone number"})
// 	return
// }

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userCtx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// call ValidateInput function here
		if errValidation := utils.ValidateInput(user, ctx); errValidation != nil {
			return
		}

		// this code is expected to work just as the previous one!
		count1, _ := checkExistingUser(userCtx, "email", *user.Email)
		count2, _ := checkExistingUser(userCtx, "phone", *user.Phone)

		if count1 > 0 || count2 > 0{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "a user with the email or phone number already exist"})
			return
		}

		password := utils.HashPassword(*user.Password)
		user.Password = &password

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		token, refreshToken, err := utils.GenerateAuthToken(*user.Email, *user.User_Role, *&user.User_id)
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while generating authentication token"})
			return
		}

		user.Token = &token
		user.Refresh_token = &refreshToken

		_, insertError := userCollection.InsertOne(userCtx, user)
		if insertError != nil {
			msg := fmt.Sprintf("unsuccessful, user item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		ctx.JSON(http.StatusOK, user)
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userCtx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user model.LoginUser
		var foundUser model.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(userCtx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "authentication failed, email or password incorrect"})
			return
		}

		validPassword, msg := utils.VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !validPassword {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := utils.GenerateAuthToken(*foundUser.Email, *foundUser.User_Role, *&foundUser.User_id)
		utils.UpdateAllToken(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"id": foundUser.ID}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, foundUser)
	}
}
