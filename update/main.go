package main

import (
	handlers "github.com/Obdurat/Schedules/handlers"
	logs "github.com/Obdurat/Schedules/logs"
	"github.com/sirupsen/logrus"
)

var file = logs.Init()

func main() {
	defer file.Close()
	handlers.Repo.Ping()
	r := handlers.Router()
	r.PUT("/schedules/:company/:id", handlers.UpdateSchedule)
	r.StaticFile("/logs", "./logs/logs.log")
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}