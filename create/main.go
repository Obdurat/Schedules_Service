package main

import (
	handlers "github.com/Obdurat/Schedules/handlers"
	"github.com/Obdurat/Schedules/logs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var _ = logs.Init()

func main() {
	ginMode := "release"
	gin.SetMode(ginMode)
	r := gin.New()
	r.Use(gin.Recovery())
	handlers.Repo.Ping()
	r.POST("/schedules/:company", handlers.CreateSchedule)
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}