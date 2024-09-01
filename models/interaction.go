package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Interaction model : Interaction/Meet related fields
type Interaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	InteractionId string             `bson:"interaction_id" json:"interaction_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CustomerID    primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Title         *string            `bson:"title,omitempty" json:"title,omitempty"`
	Description   *string            `bson:"description" json:"description"`
	StartTime     time.Time          `bson:"start_time,omitempty" json:"start_time,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
