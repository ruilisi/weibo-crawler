package controller

import (
	"go-crawler/args"
	"go-crawler/db"
	"go-crawler/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AddBloggersHandler add bloggers to database for further crawler job
func AddBloggersHandler(c *gin.Context) {
	var p args.AddBloggers

	if err := c.ShouldBind(&p); err != nil {
		StatusError(c, http.StatusBadRequest, "failed", "args bind failed")
		return
	}

	for _, blogger := range p.Bloggers {
		var bg models.Blogger
		bg.ID = blogger
		db.DB.Create(&bg)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
