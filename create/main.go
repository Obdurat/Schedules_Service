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
)

var repo *mongo.Repository

func CreateSchedule(res http.ResponseWriter, req *http.Request) {
	cn := mux.Vars(req)["company"]
	logrus.Warnf("Creating Schedule on %v", cn)
	var schedule entity.Schedule
	if err := json.NewDecoder(req.Body).Decode(&schedule); err != nil {
		logrus.Errorf("Error decoding Schedule %v: %v", cn, err)
		r := responses.ErrorRespose{Error: err.Error()}; r.Send(res, http.StatusBadRequest); return
	}
	if ve := schedule.Validate(); ve != nil {
		logrus.Errorf("Error validating Schedule on %v: %v", cn, ve)
		r := responses.ErrorRespose{Error: ve}; r.Send(res, http.StatusBadRequest); return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	col := repo.Client.Database(cn).Collection("schedules")
	if _, err := col.InsertOne(ctx, schedule); err != nil {
		logrus.Errorf("Error Inserting schedule to MongoDB: %v %v", cn, err)
		r := responses.ErrorRespose{Error: err.Error()}; r.Send(res, http.StatusBadRequest); return
	}
	defer logs.Elapsed("CreateSchedule")()
	r := responses.SuccessResponse{Message: "Created successfully"}; r.Send(res, http.StatusCreated); return
}

func main() {
	logs.Init()
	repo = mongo.New()
	logrus.Warn("Starting server...")
	repo.Ping()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{company}", CreateSchedule).Methods("POST")
	logrus.Info("Started on port 8080")
	http.ListenAndServe(":8080", router)
}