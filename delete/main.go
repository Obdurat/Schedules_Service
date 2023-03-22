package main

import (
	handlers "github.com/Obdurat/Schedules/handlers"
	logs "github.com/Obdurat/Schedules/logs"
	"github.com/sirupsen/logrus"
)

var file = logs.Init()

func main() {
	defer file.Close()
	r := handlers.Router()
	handlers.Repo.Ping()
	r.DELETE("/schedules/:company/:id", handlers.DeleteSchedule)
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}