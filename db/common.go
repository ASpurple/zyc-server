package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// const connstr = "root:xxxx@tcp(127.0.0.1:3306)/my"
const connstr = "zhou_yucheng:xxxx@tcp(127.0.0.1:3306)/zyc"

var db *sql.DB

func init() {
	d, err := sql.Open("mysql", connstr)
	if err != nil {
		panic(err)
	}
	db = d
	ReadTags()
}
