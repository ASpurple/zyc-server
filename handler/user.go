package handler

import (
	"main/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ManagerLogin 超级管理员登陆
func ManagerLogin(g *gin.Context) {
	if t := token.RootLogin(g); t != "" {
		g.JSON(http.StatusOK, succ(t))
	}
}
