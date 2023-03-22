package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	entity "github.com/Obdurat/Schedules/entity"
	logs "github.com/Obdurat/Schedules/logs"
	mongo "github.com/Obdurat/Schedules/mongo"
	responses "github.com/Obdurat/Schedules/responses"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var file = logs.Init()

var repo *mongo.Repository

func UpdateSchedule(res http.ResponseWriter, req *http.Request) {
	cn, id := mux.Vars(req)["company"], mux.Vars(req)["id"]
	logrus.Warnf("Updating Schedule on %v %v", cn, id)
	var schedule entity.Schedule
	if err := json.NewDecoder(req.Body).Decode(&schedule); err != nil {
		logrus.Errorf("Error decoding Payload on %v %v", cn, err)
		r := responses.ErrorRespose{Error: err.Error()}; r.Send(res, http.StatusBadRequest); return
	}
	if ve := schedule.Validate(); ve != nil {
		logrus.Errorf("Error validating Payload on %v %v", cn, ve)
		r := responses.ErrorRespose{Error: ve}; r.Send(res, http.StatusBadRequest); return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	obj_id, err := primitive.ObjectIDFromHex(id); if err != nil {
		logrus.Errorf("Error validating Id on %v %v: %v", cn, id, err)
		r := responses.ErrorRespose{Error: err.Error()}; r.Send(res, http.StatusBadRequest); return
	}
	col := repo.Client.Database(cn).Collection("schedules")
	e := col.FindOneAndUpdate(ctx, bson.M{"_id": obj_id}, bson.D{{Key: "$set", Value: schedule}}); if e.Err() != nil {
		logrus.Errorf("Error Updating schedule on MongoDB: %v %v: %v", cn, id, e.Err().Error())
		r := responses.ErrorRespose{Error: e.Err().Error()}; r.Send(res, http.StatusBadRequest); return
	}
	defer logrus.Info("Entry Updated on ", cn, id)
	defer logs.Elapsed("UpdateSchedule")()
	r := responses.SuccessResponse{Message: "Updated successfully"}; r.Send(res, http.StatusCreated); return
}

func main() {
	defer file.Close()
	repo = mongo.New()
	logrus.Warn("Starting server...")
	repo.Ping()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{company}/{id}", UpdateSchedule).Methods("PUT")
	logrus.Info("Started on port 8080")
	http.ListenAndServe(":8080", router)
}