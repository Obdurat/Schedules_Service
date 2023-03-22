package query

import (
	"go.mongodb.org/mongo-driver/bson"
)

func Builder(q string) (interface{}, error) {
	var where interface{}
	if len(q) == 0 {
		where = bson.M{}
		return where, nil
	}
	err := bson.UnmarshalExtJSON([]byte(q), true, &where); if err != nil {
		return nil, err
	}
	return where, nil
}