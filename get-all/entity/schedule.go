package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Client_id string `json:"client_id" bson:"client_id" validate:"required"`
	Service primitive.ObjectID `json:"service" bson:"service" validate:"required"`
	Price float32 `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	Date string `json:"date" bson:"date" validate:"required"`
	Finished bool `json:"finished" bson:"finished"`
}