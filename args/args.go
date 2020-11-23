package args

import (
	"bytes"
	"encoding/json"
	"go-crawler/utils"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Param(c *gin.Context, key string) string {
	if c.ContentType() == binding.MIMEJSON {
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		var sec map[string]string
		if err := json.Unmarshal(bodyBytes, &sec); err == nil {
			v, exist := sec[key]
			if exist {
				return v
			}
		}
	}
	return c.Request.FormValue(key)
}

func Params(c *gin.Context) map[string]string {
	sec := make(map[string]string)
	if c.ContentType() == binding.MIMEJSON {
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		if err := json.Unmarshal(bodyBytes, &sec); err == nil {
			return sec
		}
	}
	for k, v := range c.Request.URL.Query() {
		sec[k] = v[0]
	}
	return sec
}

//SetCookies args for set cookies
type SetCookies struct {
	Cookies []utils.Cookie `form:"cookies" json:"cookies" xml:"cookies" binding:"required"`
}

//QueryTags args for query tags
type QueryTags struct {
	Tags []string `form:"tags" json:"tags" xml:"tags" binding:"required"`
}

//AddBloggers args for add bloggers
type AddBloggers struct {
	Bloggers []string `form:"bloggers" json:"bloggers" xml:"bloggers" binding:"required"`
}

//QueryBlog for set parameters for query, notice that type string or []string
type QueryBlog struct {
	StartTime    string   `form:"start_time" json:"start_time" xml:"start_time"`
	EndTime      string   `form:"end_time" json:"end_time" xml:"end_time"`
	LastTime     string   `form:"last_time" json:"last_time" xml:"last_time"`
	Amount       string   `form:"amount" json:"amount" xml:"amount"`
	SuperTopic   []string `form:"super_topic" json:"super_topic" xml:"super_topic"`
	Topic        []string `form:"topic" json:"topic" xml:"topic"`
	Keywords     []string `form:"keywords" json:"keywords" xml:"keywords"`
	CategoryName string   `fomr:"category" json:"category" xml:"category"`
}

//SetCategoriesName args to rename categories
type SetCategoriesName struct {
	Names []string `form:"name" json:"names" xml:"names"`
}

//SetKeywords args for add keywords
type SetKeywords struct {
	Keywords []string `form:"keywords" json:"keywords" xml:"keywords"`
}
