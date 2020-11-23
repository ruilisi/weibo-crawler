package models

import (
	"bufio"
	"fmt"
	"go-crawler/db"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Category struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

type QueryCategory struct {
	ID    int
	Name  string
	Count int
}

type QuerySlice []QueryCategory

func GetCategories() error {
	db.DropCateGory()
	file, err := os.Open("python/dict.txt")
	defer file.Close()

	br := bufio.NewReader(file)
	cate := []string{}
	for {
		b, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		cate = append(cate, string(b))
	}

	reg, err := regexp.Compile(`\[.*?\]`)
	regID, err := regexp.Compile(`[-]?\d+`)
	regName, err := regexp.Compile(`\/[\p{Han}a-zA-Z]+`) // \p{Han} find chinese...
	for _, i := range cate {
		var category Category
		id, _ := strconv.Atoi(regID.FindString(i))
		category.ID = id + 2
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
			t, _ := FindTag(j, 2)
			t.CategoryID = category.ID
			db.DB.Save(&t)
		}
	}

	return err
}

//FetchCategories get categories and the hit blogs of last 24 hours
func FetchCategories() ([]QueryCategory, SearchResult) {
	var queryCategories []QueryCategory
	var categories []Category
	result := Result(db.DB.Order("id asc").Find(&categories).Error)
	for _, category := range categories {
		var count int64
		var tmp int64
		var tags []Tag
		db.DB.Table("tags").Where("category_id = ?", category.ID).Find(&tags)
		for _, tag := range tags {
			db.DB.Table("blogs").Where("? = ANY(tags)", tag.ID).Where("timestamp > ?", time.Now().Add(-24*time.Hour)).Count(&tmp)
			count += tmp
		}
		var qC QueryCategory
		qC.ID = category.ID
		qC.Name = category.Name
		strCount := strconv.FormatInt(count, 10)
		countInt, _ := strconv.Atoi(strCount)
		qC.Count = countInt
		queryCategories = append(queryCategories, qC)
	}
	sort.Sort(QuerySlice(queryCategories))
	return queryCategories, result
}

//SetCategories rename categories
func SetCategories(categoriesName []string) error {
	var err error

	for i, name := range categoriesName {
		tmp := db.DB
		category := Category{} //gorm feature, &category will add PK constraint to query sql
		err := tmp.Where(" id = ?", i+1).First(&category).Error
		if err != nil {
			fmt.Println("err", err)
			break
		}
		category.Name = name
		tmp.Save(&category)
	}

	if err != nil {
		return err
	}
	return nil
}

func FindCategoryByName(name string) (*Category, SearchResult) {
	var category Category
	result := Result(db.DB.Where("name = ?", name).First(&category).Error)
	return &category, result
}

func (a QuerySlice) Len() int {
	return len(a)
}
func (a QuerySlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a QuerySlice) Less(i, j int) bool {
	return a[j].Count < a[i].Count
}
