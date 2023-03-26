package handlers

import (
	"context"
	"net/http"
	"time"

	entity "github.com/Obdurat/Schedules/entity"
	logs "github.com/Obdurat/Schedules/logs"
	query "github.com/Obdurat/Schedules/query"
	"github.com/Obdurat/Schedules/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetSchedules(c *gin.Context) {
	defer logs.Elapsed("UpdateSchedule")()
	cn := c.Param("company")
	where, err := query.Builder(c.Query("where")); if err != nil {
		logrus.Errorf("Error reading query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	repo := repository.Instance
	logrus.Warnf("Reading Schedules on %v", cn)
	var schedules []entity.Schedule
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	col := repo.Database(cn)
	cursor, e := col.Find(ctx, where); if e != nil {
		logrus.Errorf("Error Reading schedules on MongoDB: %v %v", cn, e)
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()}); return
	}
	for cursor.Next(ctx) {
		var schedule entity.Schedule
		cursor.Decode(&schedule)
		schedules = append(schedules, schedule)
	}
	if err := cursor.Err(); err != nil {
		logrus.Errorf("Cursor error reading schedules on %v: %v", cn, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	defer logrus.Info("Query Complete successfully on ", cn)
	c.JSON(http.StatusOK, gin.H{"result": schedules}); return
}