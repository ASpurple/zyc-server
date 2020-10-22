package db

import (
	"database/sql"
	"fmt"
	"main/app"
	"sort"
)

// Tags ..
var Tags app.TagArr

func queryTags() (tags []app.Tag) {
	s := "SELECT * FROM tag"
	rows, err := db.Query(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		tag := app.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Level, &tag.Order, &tag.Logo, &tag.ParentID)
		if err == nil {
			tags = append(tags, tag)
		}
	}
	return
}

func setTagTree(tags app.TagArr) {
	mp := make(map[int]app.Tag)
	var arr app.TagArr
	Tags = make([]app.Tag, 0)
	for _, t := range tags {
		if t.Level == 1 {
			mp[t.ID] = t
			continue
		}
		arr = append(arr, t)
	}
	for _, t := range arr {
		if p, ok := mp[t.ParentID]; ok {
			p.Children = append(p.Children, t)
			mp[t.ParentID] = p
		}
	}
	for _, v := range mp {
		Tags = append(Tags, v)
	}
}

// ReadTags ..
func ReadTags() {
	setTagTree(queryTags())
	sort.Sort(Tags)
	for _, t := range Tags {
		sort.Sort(t.Children)
	}
}

// AddTag ..
func AddTag(tag app.Tag) error {
	s := "INSERT INTO tag(`name`,`level`,`order`,`logo`,`parentid`) VALUES(?,?,?,?,?)"
	stmt, err := db.Prepare(s)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tag.Name, tag.Level, tag.Order, tag.Logo, tag.ParentID)
	return err
}

// DelTag ..
func DelTag(id int) error {
	s := "DELETE FROM tag WHERE id = ?"
	stmt, err := db.Prepare(s)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

// UpdateTag ..
func UpdateTag(tag app.Tag) error {
	s := "UPDATE tag SET `name`=?,`level`=?,`order`=?,`logo`=?,`parentid`=? WHERE id=?"
	stmt, err := db.Prepare(s)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tag.Name, tag.Level, tag.Order, tag.Logo, tag.ParentID, tag.ID)
	return err
}

// GetTags ..
func GetTags(level int) (result []app.Tag, err error) {
	var rows *sql.Rows
	if level < 0 {
		s := "SELECT * FROM tag"
		rows, err = db.Query(s)
	} else {
		s := "SELECT * FROM tag WHERE level = ?"
		rows, err = db.Query(s, level)

	}
	if err != nil {
		return
	}
	for rows.Next() {
		t := app.Tag{}
		err := rows.Scan(&t.ID, &t.Name, &t.Level, &t.Order, &t.Logo, &t.ParentID)
		if err != nil {
			continue
		}
		result = append(result, t)
	}
	return
}

// SearchTag ..
func SearchTag(key string) (result []app.Tag, err error) {
	key = "%" + key + "%"
	s := "SELECT * FROM tag WHERE name LIKE ?"
	rows, err := db.Query(s, key)
	if err != nil {
		return
	}
	for rows.Next() {
		t := app.Tag{}
		err := rows.Scan(&t.ID, &t.Name, &t.Level, &t.Order, &t.Logo, &t.ParentID)
		if err != nil {
			continue
		}
		result = append(result, t)
	}
	return
}
