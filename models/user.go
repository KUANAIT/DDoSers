package models

import (
	"SSE/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IIN          string             `json:"iin" bson:"iin"`                     // Individual Identification Number
	IdentityCard string             `json:"identity_card" bson:"identity_card"` // Identity Card Number
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Birthday     time.Time          `json:"birthday" bson:"birthday"`
	Password     string             `json:"password" bson:"password"`
	Address      string             `json:"address" bson:"address"` // Simple address in one row
	Admin        bool               `json:"admin" bson:"admin,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// Hash sensitive data (IIN and IdentityCard)
func (u *User) HashSensitiveData() error {
	hashedIIN, err := auth.HashPassword(u.IIN)
	if err != nil {
		return err
	}
	u.IIN = hashedIIN

	hashedIdentityCard, err := auth.HashPassword(u.IdentityCard)
	if err != nil {
		return err
	}
	u.IdentityCard = hashedIdentityCard

	return nil
}

func (u *User) CheckPassword(providedPassword string) bool {
	return auth.CheckPassword(u.Password, providedPassword)
}

// Check sensitive data (for verification purposes)
func (u *User) CheckIIN(providedIIN string) bool {
	return auth.CheckPassword(u.IIN, providedIIN)
}

func (u *User) CheckIdentityCard(providedIdentityCard string) bool {
	return auth.CheckPassword(u.IdentityCard, providedIdentityCard)
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}
