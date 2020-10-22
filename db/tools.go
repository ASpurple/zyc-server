package db

import (
	"main/app"
)

// GetTools ..
func GetTools(s string, args ...interface{}) (Tools []app.Tool, err error) {
	sq := "SELECT `id`,`name`,`desc`,`link`,`logo`,`top` FROM tool LIMIT ?,?"
	if s != "" {
		sq = "SELECT `id`,`name`,`desc`,`link`,`logo`,`top` FROM tool WHERE " + s
	}
	rows, err := db.Query(sq, args...)
	for rows.Next() {
		a := app.Tool{}
		err := rows.Scan(&a.ID, &a.Name, &a.Desc, &a.Link, &a.Logo, &a.Top)
		if err == nil {
			Tools = append(Tools, a)
		}
	}
	return
}

// OperateTool 增删改
func OperateTool(s string, args ...interface{}) error {
	st, err := db.Prepare(s)
	if err != nil {
		return err
	}
	defer st.Close()
	_, err = st.Exec(args...)
	return err
}

// GetToolList 获取列表
func GetToolList(offset int, limit int) (Tools []app.Tool, err error) {
	return GetTools("", offset, limit)
}

// SearchTool 搜索
func SearchTool(words string, offset int, limit int) (Tools []app.Tool, err error) {
	words = "%" + words + "%"
	s := "CONCAT(`name`,`desc`) LIKE ? LIMIT ?,?"
	return GetTools(s, words, offset, limit)
}

// DelTool 删除
func DelTool(id int) error {
	s := "DELETE FROM tool WHERE id = ?"
	return OperateTool(s, id)
}

// AddTool 添加
func AddTool(Tool *app.Tool) error {
	s := "INSERT INTO tool(`name`,`desc`,`link`,`logo`,`top`) VALUES(?,?,?,?,?)"
	return OperateTool(s, Tool.Name, Tool.Desc, Tool.Link, Tool.Logo, Tool.Top)
}

// UpdateTool 修改
func UpdateTool(Tool *app.Tool) error {
	s := "UPDATE Tool SET `name`=?,`desc`=?,`link`=?,`logo`=?,`top`=? WHERE id = ?"
	return OperateTool(s, Tool.Name, Tool.Desc, Tool.Link, Tool.Logo, Tool.Top, Tool.ID)
}

// GetTopTool 获取所有置顶
func GetTopTool() (Tools []app.Tool, err error) {
	s := "top = 1"
	return GetTools(s)
}

// GetIndexPageTools ..
func GetIndexPageTools() (Tools []app.Tool, err error) {
	s := "top = 0"
	return GetTools(s)
}
