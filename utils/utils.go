package utils

import (
	"booking-app/database"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var JWT_SECRET_KEY string = os.Getenv("JWT_SECRET_KEY")

// Function to hashing a password
func HashPassword(password string) string {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytePassword)
}

// Function to verify or compare passwords
func VerifyPassword(existingPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(existingPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

// Function, generate authentication token based on the params passed
type SignedDetails struct {
	Email     string
	User_Role string
	ID        string
	jwt.StandardClaims
}

func GenerateAuthToken(email, role, id string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		User_Role: role,
		ID:        id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		log.Panic(err.Error())
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		log.Panic(err.Error())
		return
	}

	return token, refreshToken, err
}

func UpdateAllToken(signedToken, signedRefreshToken, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var updatedObj primitive.D

	updatedObj = append(updatedObj, bson.E{"token", signedToken})
	updatedObj = append(updatedObj, bson.E{"refresh_token", signedRefreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedObj = append(updatedObj, bson.E{"updated_at", updated_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updatedObj},
		},
		&opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

	return
}

// Input validation function
var validate = validator.New()

func ValidateInput(structField interface{}, ctx *gin.Context) error {
	validationErr := validate.Struct(structField)
	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return validationErr
	}
	return nil
}
