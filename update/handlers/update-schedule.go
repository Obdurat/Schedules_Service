package handlers

import (
	"context"
	"net/http"
	"time"

	entity "github.com/Obdurat/Schedules/entity"
	logs "github.com/Obdurat/Schedules/logs"
	"github.com/Obdurat/Schedules/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateSchedule(c *gin.Context) {
	defer logs.Elapsed("UpdateSchedule")()
	cn, id := c.Param("company"), c.Param("id")
	logrus.Warnf("Updating Schedule on %v %v", cn, id)
	var schedule entity.Schedule
	if err := c.BindJSON(&schedule); err != nil {
		logrus.Errorf("Error decoding Payload on %v %v", cn, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if ve := schedule.Validate(); ve != nil {
		logrus.Errorf("Error validating Payload on %v %v", cn, ve)
		c.JSON(http.StatusBadRequest, gin.H{"error": ve}); return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	obj_id, err := primitive.ObjectIDFromHex(id); if err != nil {
		logrus.Errorf("Error validating Id on %v %v: %v", cn, id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	col := repository.Instance.Database(cn)
	e := col.FindOneAndUpdate(ctx, bson.M{"_id": obj_id}, bson.D{{Key: "$set", Value: schedule}}); if e.Err() != nil {
		logrus.Errorf("Error Updating schedule on MongoDB: %v %v: %v", cn, id, e.Err().Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Err().Error()}); return
	}
	defer logrus.Infof("Entry Updated on %v %v", cn, id)
	c.JSON(http.StatusCreated, gin.H{"message": "Entry updated sucessfully"}); return
}