package controller

import (
	"go-crawler/args"
	"go-crawler/db"
	"go-crawler/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func SaveToRedisHandler(c *gin.Context) {
// err := models.CacheTags()
// if err != nil {
// StatusError(c, http.StatusBadRequest, "fail", err.Error())
// return
// }
//
// c.JSON(http.StatusOK, gin.H{"status": "success"})
// }

//SaveKeywordToRedisHandler cache data to redis
func SaveKeywordToRedisHandler(c *gin.Context) {
	err := models.CachKeywords()
	if err != nil {
		StatusError(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

//SetKeywordsHandler set data to database
func SetKeywordsHandler(c *gin.Context) {
	var setKeywords args.SetKeywords
	if err := c.ShouldBind(&setKeywords); err != nil {
		StatusError(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	for _, keyword := range setKeywords.Keywords {
		var tag models.Tag
		tag.Name = keyword
		tag.TagType = models.KEYWORD
		db.DB.Create(&tag)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
