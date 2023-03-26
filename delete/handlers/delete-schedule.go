package handlers

import (
	"context"
	"net/http"
	"time"

	logs "github.com/Obdurat/Schedules/logs"
	repository "github.com/Obdurat/Schedules/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteSchedule(c *gin.Context) {
	defer logs.Elapsed("DeleteSchedule")()
	repo := repository.Instance
	cn, id := c.Param("company"), c.Param("id")
	logrus.Warnf("Deleting Schedule on %v %v", cn, id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	obj_id, err := primitive.ObjectIDFromHex(id); if err != nil {
		logrus.Errorf("Error validating Id %v %v ", cn, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	col := repo.Database(cn)
	e := col.FindOneAndDelete(ctx, bson.M{"_id": obj_id}); if e.Err() != nil {
		logrus.Errorf("Error Deleting schedule on MongoDB: %v %v: %v", cn, id, e.Err().Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Err().Error()}); return
	}
	defer logrus.Infof("Entry Deleted on %v %v", cn, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted sucessfully"}); return
}