package main

import (
	handlers "github.com/Obdurat/Schedules/handlers"
	"github.com/Obdurat/Schedules/logs"
	"github.com/sirupsen/logrus"
)

var file = logs.Init()

func main() {
	defer file.Close()
	r := handlers.Router()
	handlers.Repo.Ping()
	r.POST("/schedules/:company", handlers.CreateSchedule)
	r.StaticFile("/logs", "./logs/logs.log")
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}