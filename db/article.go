package db

import (
	"fmt"
	"main/app"
)

// ArticleCount 查询文章总数
func ArticleCount() int {
	row, err := db.Query("SELECT COUNT(*) FROM Article")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	row.Next()
	c := 0
	row.Scan(&c)
	row.Close()
	return c
}

// GetArticles ..
func GetArticles(s string, args ...interface{}) (Articles []app.Article, err error) {
	sq := "SELECT `id`,`title`,`tags`,`readed`,`time`,`cover`,`top` FROM article LIMIT ?,?"
	if s != "" {
		sq = "SELECT `id`,`title`,`tags`,`readed`,`time`,`cover`,`top` FROM article WHERE " + s
	}
	rows, err := db.Query(sq, args...)
	for rows.Next() {
		a := app.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Tags, &a.Readed, &a.Time, &a.Cover, &a.Top)
		if err == nil {
			Articles = append(Articles, a)
		}
	}
	return
}

// GetArticleDetail 获取单部文章详情
func GetArticleDetail(s string, args ...interface{}) (b *app.Article, err error) {
	rows, err := db.Query(s, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	rows.Next()
	a := app.Article{}
	err = rows.Scan(&a.ID, &a.Title, &a.Tags, &a.Readed, &a.Time, &a.Cover, &a.Content, &a.Top)
	b = &a
	return
}

// OperateArticle 增删改
func OperateArticle(s string, args ...interface{}) error {
	st, err := db.Prepare(s)
	if err != nil {
		return err
	}
	defer st.Close()
	_, err = st.Exec(args...)
	return err
}

// GetArticleList 获取文章列表
func GetArticleList(offset int, limit int) (Articles []app.Article, err error) {
	return GetArticles("", offset, limit)
}

// GetTagArticle 根据标签获取文章
func GetTagArticle(tagName string, offset int, limit int) (Articles []app.Article, err error) {
	s := "FIND_IN_SET(?,tags) LIMIT ?,?"
	return GetArticles(s, tagName, offset, limit)
}

// SearchArticle 搜索
func SearchArticle(words string, offset int, limit int) (Articles []app.Article, err error) {
	words = "%" + words + "%"
	s := "CONCAT(title,tags) LIKE ? LIMIT ?,?"
	return GetArticles(s, words, offset, limit)
}

// DelArticle 删除文章
func DelArticle(id int) error {
	s := "DELETE FROM article WHERE id = ?"
	return OperateArticle(s, id)
}

// AddArticle 添加文章
func AddArticle(article *app.Article) error {
	s := "INSERT INTO article(`title`,`tags`,`readed`,`time`,`cover`,`content`,`top`) VALUES(?,?,?,?,?,?,?)"
	return OperateArticle(s, article.Title, article.Tags, article.Readed, article.Time, article.Cover, article.Content)
}

// UpdateArticle 修改文章
func UpdateArticle(article *app.Article) error {
	s := "UPDATE article SET `title`=?,`tags`=?,`readed`=?,`time`=?,`cover`=?,`content`=?,`top`=? WHERE id = ?"
	return OperateArticle(s, article.Title, article.Tags, article.Readed, article.Time, article.Cover, article.Content, article.Top, article.ID)
}

// ArticleDetail 文章详情
func ArticleDetail(id int) (Article *app.Article, err error) {
	s := "SELECT * FROM article WHERE id = ?"
	return GetArticleDetail(s, id)
}

// GetTopArticle 获取所有置顶文章
func GetTopArticle() (Articles []app.Article, err error) {
	s := "top = 1"
	return GetArticles(s)
}

// GetIndexPageArticles ..
func GetIndexPageArticles() (Articles []app.Article, err error) {
	s := "top = 0 ORDER BY id DESC"
	return GetArticles(s)
}
