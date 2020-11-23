package controller

import "github.com/gin-gonic/gin"

//StatusError an err output function of gin framework
func StatusError(c *gin.Context, httpcode int, status string, err string) {
	c.JSON(httpcode, gin.H{
		"status:": status,
		"error":   err,
	})
}
