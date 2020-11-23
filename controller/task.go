package controller

import (
	"go-crawler/tasks"
	"net/http"

	"github.com/gin-gonic/gin"
)

//TaskHandler work, work, start crawler job, only recommended in local enviroment
func TaskHandler(c *gin.Context) {
	tasks.Fetch()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
