package main

import (
	"main/handler"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := createGinEngine()
	engine.router.Use(cors)
	api(engine.router)
	engine.run()
}

type ginEngine struct {
	router *gin.Engine
	port   string
}

func (g *ginEngine) run() {
	g.router.Run(g.port)
}

func createGinEngine() *ginEngine {
	curOS := runtime.GOOS
	port := ":8080"
	if curOS != "windows" {
		// port = ":80"
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	return &ginEngine{router: router, port: port}
}

func api(router *gin.Engine) {

	router.Static("/file", "./file")

	router.POST("/upload", handler.UploadFile)

	router.POST("/login", handler.ManagerLogin)

	// level=?&keywords=? 或都不传
	router.GET("/tag", handler.GetTags)
	router.GET("/tag/tree", handler.GetTagTree)
	router.DELETE("/tag/del/:id", handler.DelTag)
	router.POST("/tag/update", handler.UpdateTag)
	router.POST("/tag/add", handler.AddTag)

	router.GET("/index", handler.IndexPageData)
	// tag=?&keywords=?&page=? 或都不传
	router.GET("/article", handler.GetArticles)
	router.GET("/article/detail/:id", handler.GetArticleDetail)
	router.POST("/article/add", handler.AddArticle)
	router.POST("/article/update", handler.UpdateArticle)
	router.DELETE("/article/del/:id", handler.DelArticle)

	router.GET("/tools", handler.ToolsPageData)
	// keywords=?&page=? 或都不传
	router.GET("/tool", handler.GetTools)
	router.POST("/tool/add", handler.AddTool)
	router.POST("/tool/update", handler.UpdateTool)
	router.DELETE("/tool/del/:id", handler.DelTool)

}

// 跨域
func cors(c *gin.Context) {
	method := c.Request.Method

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	// 处理请求
	c.Next()
}
