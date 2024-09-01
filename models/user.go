package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model :  User related fields
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Name      *string            `bson:"name" json:"name" validate:"required"`
	Password  *string            `bson:"password" json:"password" validate:"required,min=2,max=100"`
	Email     *string            `bson:"email" json:"email" validate:"email,required"`
	Role      *string            `bson:"role" json:"role" validate:"required,eq=ADMIN|eq=USER"`
	Company   *string            `bson:"company,omitempty" json:"company,omitempty"`
	PhoneNo   *string            `bson:"phone_no,omitempty" json:"phone_no,omitempty"`
	Token     *string            `bson:"token,omitempty" json:"token,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
