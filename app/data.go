package app

import "html/template"

// Tag ..
type Tag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Order    int    `json:"order"`
	Logo     string `json:"logo"`
	ParentID int    `json:"parentid"`
	Children TagArr `json:"children"`
}

// TagArr ..
type TagArr []Tag

func (s TagArr) Len() int {
	return len(s)
}

func (s TagArr) Less(i, j int) bool {
	return s[i].Order < s[j].Order
}

func (s TagArr) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Article ..
type Article struct {
	ID      int           `json:"id"`
	Title   string        `json:"title"`
	Tags    string        `json:"tags"`
	Readed  int           `json:"readed"`
	Time    string        `json:"time"`
	Cover   string        `json:"cover"`
	Content template.HTML `json:"content"`
	Top     int           `json:"top"`
}

// Tool ..
type Tool struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Link string `json:"link"`
	Logo string `json:"logo"`
	Top  int    `json:"top"`
}

// SuccRes 请求成功返回数据结构
type SuccRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ErrRes 请求失败返回数据结构
type ErrRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
