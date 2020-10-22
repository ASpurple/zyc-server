package handler

import (
	"main/app"
	"main/db"
	"main/token"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var articlePageSize = 12

// GetTopArticle 获取所有置顶文章
func GetTopArticle(c *gin.Context) {
	list, err := db.GetTopArticle()
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ(list))
}

// GetArticles ..
func GetArticles(c *gin.Context) {
	keywords := c.Query("keywords")
	tag := c.Query("tag")
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	var list []app.Article
	if tag != "" {
		list, err = db.GetTagArticle(tag, (page-1)*articlePageSize, articlePageSize)
	} else if keywords != "" {
		list, err = db.SearchArticle(keywords, (page-1)*articlePageSize, articlePageSize)
	} else {
		list, err = db.GetArticleList((page-1)*articlePageSize, articlePageSize)
	}
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ(list))
}

// GetArticleDetail handler
func GetArticleDetail(c *gin.Context) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	art, err := db.ArticleDetail(id)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ(art))
}

// AddArticle 添加
func AddArticle(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	art := app.Article{}
	err := c.BindJSON(&art)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	art.Time = time.Now().Format("2006-01-02")
	err = db.AddArticle(&art)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("添加成功"))
}

// UpdateArticle 更新
func UpdateArticle(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	art := app.Article{}
	err := c.BindJSON(&art)
	if err != nil || art.ID == 0 {
		c.JSON(http.StatusOK, failed("参数错误"))
		return
	}
	err = db.UpdateArticle(&art)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("更新成功"))
}

// DelArticle 删除
func DelArticle(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusOK, failed("参数错误"))
		return
	}
	err = db.DelArticle(id)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("删除成功"))
}

// IndexPageData ..
func IndexPageData(c *gin.Context) {
	result := make([]app.Article, 0)
	t, err := db.GetTopArticle()
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	result = t
	n, err := db.GetIndexPageArticles()
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	result = append(result, n...)
	c.JSON(http.StatusOK, succ(result))
}
