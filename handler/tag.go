package handler

import (
	"main/app"
	"main/db"
	"main/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddTag ..
func AddTag(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	t := app.Tag{}
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	err = db.AddTag(t)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("添加成功"))
}

// UpdateTag 更新
func UpdateTag(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	t := app.Tag{}
	err := c.BindJSON(&t)
	if err != nil || t.Name == "" {
		c.JSON(http.StatusOK, failed("参数错误"))
		return
	}
	err = db.UpdateTag(t)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("更新成功"))
}

// DelTag 删除
func DelTag(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	err = db.DelTag(id)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("删除成功"))
}

// GetTags handler
func GetTags(c *gin.Context) {
	words := c.Query("keywords")
	level := c.Query("level")
	var ts []app.Tag
	le, err := strconv.Atoi(level)
	if err != nil {
		le = -1
	}
	if words == "" {
		ts, err = db.GetTags(le)
	} else {
		ts, err = db.SearchTag(words)
	}
	c.JSON(200, succ(ts))
}

// GetTagTree ..
func GetTagTree(c *gin.Context) {
	c.JSON(http.StatusOK, succ(db.Tags))
}
