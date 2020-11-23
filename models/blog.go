package models

import (
	"go-crawler/db"
	"time"

	"github.com/lib/pq"
)

type Blog struct {
	ID        string `gorm:"primaryKey"`
	BloggerID string
	Tags      pq.StringArray `gorm:"type:varchar(255)[]"`
	Content   string
	Timestamp time.Time
	Pictures  pq.StringArray `gorm:"type:varchar(255)[]"`
	Video     string
	Article   string
	Rt        string
}

type BlogQuery struct {
	BlogID    string
	BloggerID string
	Username  string
	Follows   int64
	Avatar    string
	Tags      pq.StringArray
	Content   string
	Timestamp time.Time
	Pictures  pq.StringArray
	Video     string
	Article   string
	Rt        string
}

func FindBlogByID(ID interface{}) (*Blog, SearchResult) {
	var blog Blog
	result := Result(db.DB.Where("id = ?", ID).First(&blog).Error)
	return &blog, result
}

func FindBlogsByBloggerID(bgid string) ([]Blog, SearchResult) {
	var blogs []Blog
	result := Result(db.DB.Where("bloggerid = ?", bgid).Find(&blogs).Error)
	return blogs, result
}

func FindBlogsByTags(tags []string) ([]Blog, SearchResult) {
	var blogs []Blog
	tmpdb := db.DB
	for _, tag := range tags {
		tmpdb = tmpdb.Where("? = ANY(tags)", tag)
	}
	result := Result(tmpdb.Find(&blogs).Error)
	return blogs, result
}

//FindBlogs due to unknown reason of lose of string array query data, do query twice to fill the lost
func FindBlogs(st time.Time, et time.Time, cnt int, supertopic []string, topic []string, keywords []string, categoryName string) ([]BlogQuery, SearchResult) {
	var blogs []BlogQuery
	var bs []Blog
	var tags []Tag
	var tagsName []Tag
	var category *Category

	tagsName, _ = GetTags()

	// get tags query parameters
	if !(len(supertopic)+len(topic)+len(keywords) == 0 && categoryName == "") {
		tags = append(tags, addTags(supertopic, 0)...)
		tags = append(tags, addTags(topic, 1)...)
		tags = append(tags, addTags(keywords, 2)...)
		if categoryName != "" && len(keywords) == 0 {
			category, _ = FindCategoryByName(categoryName)
			tmpTags, _ := FindTagsByCategoryID(category.ID)
			tags = append(tags, tmpTags...)
		}
	}

	tmpPointer := db.DB
	_tmpPointer := db.DB

	tmpPointer = tmpPointer.Table("blogs").Select("blogs.id as blog_id,blogs.blogger_id,bloggers.username,bloggers.follows,bloggers.avatar,tags,content,timestamp,pictures,video,article,rt").Joins("join bloggers on bloggers.id = blogs.blogger_id")

	if !time.Time.IsZero(st) && !time.Time.IsZero(et) {
		tmpPointer = tmpPointer.Where("timestamp > ? AND timestamp < ?", st, et)
		_tmpPointer = _tmpPointer.Where("timestamp > ? AND timestamp < ?", st, et)
	} else if time.Time.IsZero(st) && !time.Time.IsZero(et) {
		tmpPointer = tmpPointer.Where("timestamp < ?", et)
		_tmpPointer = _tmpPointer.Where("timestamp < ?", et)
	} else if time.Time.IsZero(et) && !time.Time.IsZero(st) {
		tmpPointer = tmpPointer.Where("timestamp > ?", st)
		_tmpPointer = _tmpPointer.Where("timestamp > ?", st)
	}

	//Where . AND( . OR . OR . OR ...)
	//add tags parameters
	if len(tags) != 0 {
		tmpP := db.DB
		_tmpP := db.DB
		for i, tag := range tags {
			if i == 0 {
				tmpP = tmpP.Where("? = ANY(tags)", tag.ID.String())
				_tmpP = _tmpP.Where("? = ANY(tags)", tag.ID.String())
			} else {
				tmpP = tmpP.Or("? = ANY(tags)", tag.ID.String())
				_tmpP = _tmpP.Or("? = ANY(tags)", tag.ID.String())
			}
		}
		tmpPointer = tmpPointer.Where(tmpP)
		_tmpPointer = _tmpPointer.Where(_tmpP)
	}

	result := Result(tmpPointer.Order("timestamp desc").Limit(cnt).Find(&blogs).Error)
	_tmpPointer.Order("timestamp desc").Limit(cnt).Find(&bs)

	//add tags name
	for i := range blogs {
		blogs[i].Tags = bs[i].Tags
		rsTags := []string{}
		for _, t := range blogs[i].Tags {
			for _, tag := range tagsName {
				if t == tag.ID.String() {
					rsTags = append(rsTags, tag.Name)
				}
			}
		}

		//fill hte lost mentioned above
		blogs[i].Tags = rsTags
		blogs[i].Pictures = bs[i].Pictures
	}
	return blogs, result
}

func FindLastest(bloggerID interface{}) (*Blog, SearchResult) {
	var blog Blog
	result := Result(db.DB.Where("blogger_id = ?", bloggerID).Order("timestamp desc").First(&blog).Error)
	return &blog, result
}

func addTags(input []string, t int) []Tag {
	var pTag *Tag
	var tags []Tag
	if len(input) != 0 {
		for _, tag := range input {
			switch t {
			case 0:
				pTag, _ = FindTag(""+tag, t) //the label of super topic
			case 1:
				pTag, _ = FindTag(tag, t)
			case 2:
				pTag, _ = FindTag(tag, t)
			}
			if pTag != nil {
				tags = append(tags, *pTag)
			}
		}
	}
	return tags
}
