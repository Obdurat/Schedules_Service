package main

import (
	"context"
	"net/http"
	"time"

	entity "github.com/Obdurat/Schedules/entity"
	logs "github.com/Obdurat/Schedules/logs"
	mongo "github.com/Obdurat/Schedules/mongo"
	query "github.com/Obdurat/Schedules/query"
	responses "github.com/Obdurat/Schedules/responses"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var file = logs.Init()

var repo *mongo.Repository

func GetSchedules(res http.ResponseWriter, req *http.Request) {
	cn := mux.Vars(req)["company"]
	where, err := query.Builder(req.FormValue("where")); if err != nil {
		logrus.Errorf("Error reading query: %v", err)
		r := responses.ErrorRespose{Error: err.Error()}; r.Send(res, http.StatusBadRequest); return
	}
	logrus.Warnf("Reading Schedules on %v", cn)
	var schedules []entity.Schedule
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	col := repo.Client.Database(cn).Collection("schedules")
	cursor, e := col.Find(ctx, where); if e != nil {
		logrus.Errorf("Error Reading schedules on MongoDB: %v %v", cn, e)
		r := responses.ErrorRespose{Error: e.Error()}; r.Send(res, http.StatusBadRequest); return
	}
	for cursor.Next(ctx) {
		var schedule entity.Schedule
		cursor.Decode(&schedule)
		schedules = append(schedules, schedule)
	}
	if err := cursor.Err(); err != nil {
		logrus.Errorf("Cursor error reading schedules on %v: %v", cn, err.Error())
		r := responses.ErrorRespose{Error: e.Error()}; r.Send(res, http.StatusInternalServerError); return
	}
	defer logrus.Info("Query Complete successfully on ", cn)
	defer logs.Elapsed("UpdateSchedule")()
	r := responses.SuccessResponse{Message: schedules}; r.Send(res, http.StatusCreated); return
}

func main() {
	defer file.Close()
	repo = mongo.New()
	logrus.Warn("Starting server...")
	repo.Ping()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{company}", GetSchedules).Queries("where", "{where}").Methods("GET")
	router.HandleFunc("/schedules/{company}", GetSchedules).Methods("GET")
	logrus.Info("Started on port 8080")
	http.ListenAndServe(":8080", router)
}