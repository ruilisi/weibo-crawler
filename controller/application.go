package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//PingHandler test connection
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
