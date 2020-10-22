package handler

import (
	"main/app"
	"main/db"
	"main/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var toolPageSize = 24

// GetTopTool 获取所有置顶
func GetTopTool(c *gin.Context) {
	list, err := db.GetTopTool()
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ(list))
}

// GetTools ..
func GetTools(c *gin.Context) {
	keywords := c.Query("keywords")
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	var list []app.Tool
	if keywords != "" {
		list, err = db.SearchTool(keywords, (page-1)*toolPageSize, toolPageSize)
	} else {
		list, err = db.GetToolList((page-1)*toolPageSize, toolPageSize)
	}
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ(list))
}

// AddTool 添加
func AddTool(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	tool := app.Tool{}
	err := c.BindJSON(&tool)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	err = db.AddTool(&tool)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("添加成功"))
}

// UpdateTool 更新
func UpdateTool(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	tool := app.Tool{}
	err := c.BindJSON(&tool)
	if err != nil || tool.ID == 0 {
		c.JSON(http.StatusOK, failed("参数错误"))
		return
	}
	err = db.UpdateTool(&tool)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("更新成功"))
}

// DelTool 删除
func DelTool(c *gin.Context) {
	if !token.Root(c) {
		return
	}
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusOK, failed("参数错误"))
		return
	}
	err = db.DelTool(id)
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, succ("删除成功"))
}

// ToolsPageData ..
func ToolsPageData(c *gin.Context) {
	result := make([]app.Tool, 0)
	t, err := db.GetTopTool()
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	result = t
	n, err := db.GetIndexPageTools()
	if err != nil {
		c.JSON(http.StatusOK, failed(err.Error()))
		return
	}
	result = append(result, n...)
	c.JSON(http.StatusOK, succ(result))
}
