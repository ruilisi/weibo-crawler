package models

import (
	"fmt"
	"go-crawler/db"
	"go-crawler/redis"
)

const (
	SUPER_TOPIC = 0
	TOPIC       = 1
	KEYWORD     = 2
)

type Tag struct {
	Model
	Name       string `gorm:"unique"`
	TagType    int
	CategoryID int
}

func GetTags() ([]Tag, SearchResult) {
	var tags []Tag
	result := Result(db.DB.Find(&tags).Error)
	return tags, result
}

func FindTag(name string, t int) (*Tag, SearchResult) {
	var tag Tag
	result := Result(db.DB.Where("name =  ? AND tag_type = ?", name, t).First(&tag).Error)
	return &tag, result
}

func FindTagByID(id string) (*Tag, SearchResult) {
	var tag Tag
	result := Result(db.DB.Where("id = ?", id).First(&tag).Error)
	return &tag, result
}

func FindTagsByCategoryID(cID int) ([]Tag, SearchResult) {
	var tags []Tag
	result := Result(db.DB.Where("category_id = ?", cID).Find(&tags).Error)
	return tags, result
}

func CacheTags() error {
	tags := []Tag{}

	err := db.DB.Find(&tags).Error
	if err != nil {
		return err
	}
	for _, tag := range tags {
		redis.HSet(fmt.Sprintf("tag:%v", tag.TagType), tag.Name, tag.ID)
	}
	return nil
}

//FindTags find tags for cluster
func FindTags() ([]Tag, SearchResult) {
	var tags []Tag
	//	var rs []Tag
	//	result := Result(db.DB.Where("tag_type = ? AND tag_type = ? AND tag_name ", 1, 2).Find(&tags).Error)
	result := Result(db.DB.Where("tag_type = ?", 2).Find(&tags).Error)
	// for _, i := range tags {
	// if !(i.TagType == 1 && len(i.Name) >= 21) {
	// rs = append(rs, i)
	// }
	// }
	return tags, result
}

func CachKeywords() error {
	tags := []Tag{}

	err := db.DB.Where("tag_type = ?", KEYWORD).Find(&tags).Error
	if err != nil {
		return err
	}
	redis.Del("keywords")
	for _, tag := range tags {
		redis.SAdd("keywords", tag.Name)
	}
	return nil
}

func (t *Tag) CacheTag() {
	redis.HSet("tag:"+fmt.Sprint(t.TagType), t.Name, t.ID)
}

func CreateTag(name string, t int) (*Tag, SearchResult) {
	var tag Tag
	var res SearchResult
	foundTag, result := FindTag(name, t)
	if result.NotFound() {
		tag.Name = name
		tag.TagType = t
		createErr := db.DB.Create(&tag).Error
		if createErr == nil {
			tag.CacheTag()
		}
		res = Result(createErr)
		return &tag, res
	} else {
		return foundTag, result
	}
}
