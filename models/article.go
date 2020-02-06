package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	// 公共继承，类似于OO中的父类对象
	Model
	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`
	// 自定义字段
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreatedOn", time.Now().Unix())
	err = scope.SetColumn("ModifiedOn", time.Now().Unix())
	err = scope.SetColumn("ModifiedBy", article.CreatedBy)
	return err
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedOn", time.Now().Unix())
	err = scope.SetColumn("ModifiedBy", article.ModifiedBy)
	return err
}

func (article *Article) BeforeDelete(scope *gorm.Scope) error {
	err := scope.SetColumn("DeletedOn", time.Now().Unix())
	return err
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//db.Preload("Tag").Where("id = ?", id).First(&article)
	return &article, nil
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (article []*Article, err error) {
	err = db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(Article{}).Where(maps).Count(&count)
	return
}

func ExistsArticleById(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		CreatedBy: data["created_by"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		State:     data["state"].(int),
	})
	return true
}

func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)
	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}
