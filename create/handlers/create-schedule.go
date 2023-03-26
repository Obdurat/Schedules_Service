package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Obdurat/Schedules/entity"
	"github.com/Obdurat/Schedules/logs"
	"github.com/Obdurat/Schedules/mongo"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)



func CreateSchedule(c *gin.Context) {
	repo := mongo.Repo
	defer logs.Elapsed("CreateSchedule")()
	cn := c.Param("company")
	logrus.Warnf("Creating Schedule on %v", cn)
	var schedule entity.Schedule
	if err := c.BindJSON(&schedule); err != nil {
		logrus.Errorf("Error decoding Schedule %v: %v", cn, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := schedule.Validate(); err != nil {
		logrus.Errorf("Error validating Schedule on %v: %v", cn, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err}); return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	col := repo.Collection(cn)
	if _, err := col.InsertOne(ctx, schedule); err != nil {
		logrus.Errorf("Error Inserting schedule to MongoDB: %v %v", cn, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created sucessfully"}); return
}