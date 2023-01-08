package model

import (
	"time"
)

type User struct {
	ID            string    `bson:"id"`
	First_name    *string   `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string   `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string   `json:"password" validate:"required,min=6"`
	Email         *string   `json:"email" validate:"email,required"`
	Phone         *string   `json:"phone" validate:"required"`
	Refresh_token *string   `json:"refresh_token"`
	Token         *string   `json:"token"`
	User_Role     *string   `json:"user_role" validate:"required,eq=ADMIN|eq=USER"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

type LoginUser struct {
	Email    *string `json:"email" validate:"email,required"`
	Password *string `json:"password" validate:"required,min=6"`
}
