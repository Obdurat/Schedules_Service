package handlers

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	ginMode := "release"
	gin.SetMode(ginMode)
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}