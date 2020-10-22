package handler

import (
	"crypto/md5"
	"fmt"
	"main/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取当前日期字符串,返回年月，日
func curDateStr() (string, string) {
	y := time.Now().Year()
	m := int(time.Now().Month())
	d := time.Now().Day()
	ys := strconv.Itoa(y)
	ms := strconv.Itoa(m)
	ds := strconv.Itoa(d)
	if len(ms) == 1 {
		ms = "0" + ms
	}
	if len(ds) == 1 {
		ds = "0" + ds
	}
	return ys + ms, ds
}

// Suffix 取后缀名
func Suffix(n string) string {
	i := strings.LastIndex(n, ".")
	if i < 0 {
		return ""
	}
	return n[i:]
}

func md(t string) string {
	bs := []byte(t)
	return fmt.Sprintf("%x", md5.Sum(bs))
}

// UploadFile 上传文件
func UploadFile(g *gin.Context) {
	if !token.Root(g) {
		return
	}
	file, err := g.FormFile("file")
	if err != nil {
		g.JSON(200, failed(err.Error()))
		return
	}
	ym, d := curDateStr()
	curDir, err := os.Getwd()
	if err != nil {
		g.JSON(200, failed(err.Error()))
		return
	}
	p := filepath.Join(filepath.Join(curDir, "file"), ym, d)
	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		g.JSON(200, failed(err.Error()))
		return
	}
	name := file.Filename
	name = strings.ReplaceAll(name, "|", "")
	p = filepath.Join(p, name)
	err = g.SaveUploadedFile(file, p)
	if err != nil {
		g.JSON(200, failed(err.Error()))
		return
	}
	localPath := "/file/" + ym + "/" + d + "/" + name // 本地存储路径
	g.JSON(200, succ(localPath))
}
