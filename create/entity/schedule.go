package entity

import (
	"github.com/go-playground/validator/v10"
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

type VE struct {
	Param string `json:"param"`
	Message string `json:"message"`
}

func vError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	default:
		return fe.Field() + " is required"
	}
}

func (s *Schedule) Validate() []VE {
	v := validator.New()
	if err := v.Struct(s); err != nil {
		out := make([]VE, len(err.(validator.ValidationErrors)))
		for i, fe := range err.(validator.ValidationErrors) {
			out[i] = VE{fe.Field(), vError(fe)}
		}
		return out
	}
	return nil
}