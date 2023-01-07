package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

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
			ExpiresAt: time.Now().Local().Add(time.Hour *  time.Duration(168)).Unix(),
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
