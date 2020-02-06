package cache_service

import (
	"encoding/json"
	"gin-blog/gredis"
	"gin-blog/models"
	"gin-blog/pkg/exception"
	"gin-blog/pkg/mylog"
	"strconv"
	"strings"
)

type Article struct {
	ID    int
	TagID int
	State int

	PageNum  int
	PageSize int
}

func (a *Article) GetArticleKey() string {
	return exception.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

func (a *Article) GetArticlesKey() string {
	keys := []string{
		exception.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}
	return strings.Join(keys, "_")
}

func (a *Article) ExistsById() (bool, error) {
	return models.ExistsArticleById(a.ID)
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := Article{ID: a.ID}
	key := cache.GetArticleKey()

	// 如果有缓存，我们直接查询缓存返回结果
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			mylog.Info(err)
		} else {
			_ = json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	// 如果无缓存，我们将mysql中查询的数据结果缓存起来，再返回上层调用函数
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	_, _ = gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		article, cacheArticles []*models.Article
	)

	cache := Article{
		TagID:    a.TagID,
		State:    a.State,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()

	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			mylog.Info(err)
		} else {
			_ = json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	article, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	_, _ = gredis.Set(key, article, 3600)
	return article, nil

}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if a.State != -1 {
		maps["state"] = a.State
	}

	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}
