package controller

import (
	"go-crawler/args"
	"go-crawler/db"
	"go-crawler/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Set categories from local file which exported
// func SetCategoriesHandler(c *gin.Context) {
// err := models.GetCategories()
// if err != nil {
// StatusError(c, http.StatusBadRequest, "fail", err.Error())
// }
//
// c.JSON(http.StatusOK, gin.H{
// "status": "success",
// })
// }

//SetCategoriesHandler set categories from request
func SetCategoriesHandler(c *gin.Context) {
	db.DropCateGory()
	params := args.Params(c)
	categoryString := params["category"]
	cate := strings.Split(categoryString, "\n")

	reg, err := regexp.Compile(`\[.*?\]`)
	regID, err := regexp.Compile(`[-]?\d+`)
	regName, err := regexp.Compile(`\/[\p{Han}a-zA-Z]+`) // \p{Han} find chinese ...
	if err != nil {
		StatusError(c, http.StatusBadRequest, "fail", err.Error())
		return
	}
	for _, i := range cate {
		var category models.Category
		id, _ := strconv.Atoi(regID.FindString(i))
		category.ID = id + 2 // id starts from -1, which means noise data
		tags := reg.FindString(i)
		tags = strings.Replace(tags[1:len(tags)-1], `'`, "", -1)
		tags = strings.Replace(tags, " ", "", -1)
		tagName := strings.Split(tags, ",")
		category.Name = regName.FindString(i)
		if category.Name != "" {
			category.Name = category.Name[1:]
		}
		db.DB.Create(&category)
		for _, j := range tagName {
			t, _ := models.FindTag(j, 2)
			t.CategoryID = category.ID
			db.DB.Save(&t)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

//SetCategoriesNameHandler rename the categories
func SetCategoriesNameHandler(c *gin.Context) {
	var categoriesName args.SetCategoriesName
	if err := c.ShouldBind(&categoriesName); err != nil {
		StatusError(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	err := models.SetCategories(categoriesName.Names)

	if err != nil {
		StatusError(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
