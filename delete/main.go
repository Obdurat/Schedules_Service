package main

import (
	"context"
	"net/http"
	"time"

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

func DeleteSchedule(res http.ResponseWriter, req *http.Request) {
	cn, id := mux.Vars(req)["company"], mux.Vars(req)["id"]
	logrus.Warnf("Deleting Schedule on %v %v", cn, id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	obj_id, err := primitive.ObjectIDFromHex(id); if err != nil {
		logrus.Errorf("Error validating Id %v %v ", cn, err)
		r := responses.ErrorRespose{Error: err.Error()}; r.Send(res, http.StatusBadRequest); return
	}
	col := repo.Client.Database(cn).Collection("schedules")
	e := col.FindOneAndDelete(ctx, bson.M{"_id": obj_id}); if e.Err() != nil {
		logrus.Errorf("Error Deleting schedule on MongoDB: %v %v: %v", cn, id, e.Err().Error())
		r := responses.ErrorRespose{Error: e.Err().Error()}; r.Send(res, http.StatusBadRequest); return
	}
	defer logrus.Infof("Entry Deleted on %v %v", cn, id)
	defer logs.Elapsed("UpdateSchedule")()
	r := responses.SuccessResponse{Message: "Deleted successfully"}; r.Send(res, http.StatusOK); return
}

func main() {
	defer file.Close()
	repo = mongo.New()
	logrus.Warn("Starting server...")
	repo.Ping()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{company}/{id}", DeleteSchedule).Methods("DELETE")
	logrus.Info("Started on port 8080")
	http.ListenAndServe(":8080", router)
}