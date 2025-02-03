package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MobileNumber string             `bson:"mobile_number" json:"mobile_number" validate:"required,e164"`
	CountryCode  string             `bson:"country_code" json:"country_code" validate:"required"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	IsVerified   bool               `bson:"is_verified" json:"is_verified"`
	Roles        []string           `bson:"roles" json:"roles"`
	LastLoginAt  time.Time          `bson:"last_login_at" json:"last_login_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserRegistration struct {
	MobileNumber string `json:"mobile_number" validate:"required,e164"`
	CountryCode  string `json:"country_code" validate:"required"`
	Password     string `json:"password" validate:"required,min=8,max=72"`
	OTPCode      string `json:"otp_code" validate:"required"`
}

type UserLogin struct {
	MobileNumber string `json:"mobile_number" validate:"required,e164"`
	Password     string `json:"password" validate:"required"`
}