package models

import (
	"time"

	"github.com/josephmakin/monerohub/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Donation struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Address string `json:"address" bson:"address"`
	Amount uint64 `json:"amount" bson:"amount"`
	Message string `json:"message" bson:"message"`
	Name string `json:"name" bson:"name"`
	Paid bool `json:"paid" bson:"paid"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Transactions []models.Transaction `json:"transactions" bson:"transactions"`
}
