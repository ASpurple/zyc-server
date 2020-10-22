package token

import (
	"crypto/md5"
	"fmt"
	"main/app"
	"time"

	"github.com/gin-gonic/gin"
)

var managerName = "zhou_yucheng"
var managerKey = "0op-0op-"

var managerToken = ""

type managerInfo struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Root 超级管理员鉴权
func Root(g *gin.Context) bool {
	token := g.Request.Header.Get("Token")
	if managerToken == "" || token != managerToken {
		errRes := app.ErrRes{Code: 400, Msg: "Permission denied"}
		g.JSON(200, errRes)
		return false
	}
	return true
}

// 登陆后设置 managerToken 的值
func setManagerToken() {
	t := time.Now().UTC()
	v := t.String()
	v = v + managerName
	v = fmt.Sprintf("%x", md5.Sum([]byte(v)))
	managerToken = v
}

// RootLogin 超级管理员登陆 登陆成功返回managerToken，失败返回空字符串
func RootLogin(g *gin.Context) string {
	info := managerInfo{}
	err := g.BindJSON(&info)
	if err != nil {
		errRes := app.ErrRes{Code: -1, Msg: "Parameter error"}
		g.JSON(200, errRes)
		return ""
	}
	key := fmt.Sprintf("%x", md5.Sum([]byte(managerKey)))
	if info.Name != managerName || info.Key != key {
		errRes := app.ErrRes{Code: -1, Msg: "用户名或密码错误"}
		g.JSON(200, errRes)
		return ""
	}
	setManagerToken()
	return managerToken
}
