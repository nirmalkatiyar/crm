package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ticket model : Ticket related fields
type Ticket struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TicketId      string             `bson:"ticket_id" json:"ticket_id"`
	InteractionID primitive.ObjectID `bson:"interaction_id" json:"interaction_id"`
	CustomerID    primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Status        *string            `bson:"status" json:"status" validate:"required,eq=open|eq=in_progress|eq=resolved|eq=closed"`
	Description   *string            `bson:"description" json:"description"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
