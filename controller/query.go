package controller

import (
	"go-crawler/args"
	"go-crawler/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//QueryBlogsHandler get blogs from database with parameters
func QueryBlogsHandler(c *gin.Context) {
	var queryBlog args.QueryBlog
	if err := c.ShouldBind(&queryBlog); err != nil {
		StatusError(c, http.StatusBadRequest, "failed", err.Error())
	}

	st, _ := time.ParseInLocation("2006-01-02 15:04:05", queryBlog.StartTime, time.Local)
	et, _ := time.ParseInLocation("2006-01-02 15:04:05", queryBlog.EndTime, time.Local)
	lt, _ := time.ParseInLocation("2006-01-02 15:04:05", queryBlog.LastTime, time.Local)
	cnt, _ := strconv.Atoi(queryBlog.Amount)

	if st.Unix() == 0 {
		st = time.Time{}
	}
	if et.Unix() == 0 {
		et = time.Time{}
	}

	// if last time < end time, replace et with lt
	if et.After(lt) && !lt.IsZero() || et.IsZero() {
		et = lt
	}

	blogs, rs := models.FindBlogs(st, et, cnt, queryBlog.SuperTopic, queryBlog.Topic, queryBlog.Keywords, queryBlog.CategoryName)
	if rs.Err() != nil {
		StatusError(c, http.StatusBadRequest, "failed", rs.Err().Error())
		return
	} else if rs.NotFound() {
		StatusError(c, http.StatusBadRequest, "failed", "not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"blogs":  blogs,
	})
}

//QueryTagsHandler get all tags from database
func QueryTagsHandler(c *gin.Context) {

	tags, rs := models.GetTags()
	if rs.Err() != nil {
		StatusError(c, http.StatusBadRequest, "failed", rs.Err().Error())
		return
	} else if rs.NotFound() {
		StatusError(c, http.StatusBadRequest, "failed", "not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"tags":   tags,
	})

}

//QueryKeywordsHandler get cluster input data
func QueryKeywordsHandler(c *gin.Context) {
	tags, rs := models.FindTags()
	if rs.Err() != nil {
		StatusError(c, http.StatusBadRequest, "failed", rs.Err().Error())
		return
	} else if rs.NotFound() {
		StatusError(c, http.StatusBadRequest, "failed", "not found")
		return
	}

	file, _ := os.Create("/python/keywords.txt")
	defer file.Close()
	for _, tag := range tags {
		file.WriteString(tag.Name + "\n")
	}
	file.Sync()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"tags":   tags,
	})
}

//QueryCategroyHandler get categories from database including the hit blogs of last 24 hours
func QueryCategroyHandler(c *gin.Context) {
	categories, rs := models.FetchCategories()

	if rs.Err() != nil {
		StatusError(c, http.StatusBadRequest, "failed", rs.Err().Error())
		return
	} else if rs.NotFound() {
		StatusError(c, http.StatusBadRequest, "failed", "not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"categories": categories,
	})
}
