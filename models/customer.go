package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer model : Customer related fields
type Customer struct {
	CustomerId string             `bson:"customer_id" json:"customer_id"`
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       *string            `bson:"name" json:"name" validate:"required"`
	Email      *string            `bson:"email" json:"email" validate:"email,required"`
	Password   *string            `bson:"password" json:"password" validate:"required,min=2,max=100"`
	Company    *string            `bson:"company,omitempty" json:"company,omitempty"`
	Phone      *string            `bson:"phone,omitempty" json:"phone,omitempty"`
	Token      *string            `bson:"token,omitempty" json:"token,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
