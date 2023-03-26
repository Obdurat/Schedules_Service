package main

import (
	handlers "github.com/Obdurat/Schedules/handlers"
	logs "github.com/Obdurat/Schedules/logs"
	repository "github.com/Obdurat/Schedules/repository"
	"github.com/sirupsen/logrus"
)

var file = logs.Init()

func main() {
	logrus.Warnf("Service Starting...")
	defer file.Close()
	r := handlers.Router()
	repository.Instance.Ping()
	defer repository.Instance.Close()
	r.POST("/schedules/:company", handlers.CreateSchedule)
	r.StaticFile("/logs", "./logs/logs.log")
	logrus.Info("Service Ready")
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}