package models

import "go-crawler/db"

type Blogger struct {
	ID       string `gorm:"primaryKey"`
	Username string
	Avatar   string
	Follows  int64
}

func FindBloggerByID(id string) (*Blogger, SearchResult) {
	var blogger Blogger
	result := Result(db.DB.Where("id = ?", id).First(&blogger).Error)
	return &blogger, result
}

// func FindBloggerByName(name string) (*Blogger, SearchResult) {
// var blogger Blogger
// result := Result(db.DB.Where("name = ?", name).First(&blogger).Error)
// return &blogger, result
// }

func FindBloggers() ([]Blogger, SearchResult) {
	var blogger []Blogger
	result := Result(db.DB.Where("id != ''").Find(&blogger).Error)
	return blogger, result
}
