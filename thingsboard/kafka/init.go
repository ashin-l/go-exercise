package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var tbdb *sql.DB
var st *sql.Stmt

const CONNDB = "postgres://postgres:111@192.168.152.44/tbstress?sslmode=disable"

func InitDB() (err error) {
	tbdb, err = sql.Open("postgres", CONNDB)
	if err != nil {
		fmt.Println("connect database error:", err)
	}
	sqlstr := "insert into dvdata(deviceid, value, other, clienttime, servertime) values($1, $2, $3, $4, $5)"
	st, err = tbdb.Prepare(sqlstr)
	if err != nil {
		fmt.Println("create insert stmt error:", err)
	}
	return
}

func Syncdb() {
	createTable()
}

func createTable() {
	db, err := sql.Open("postgres", CONNDB)
	if err != nil {
		fmt.Println("database open error:", err)
		return
	}
	defer db.Close()
	_, err = db.Query("create table dvdata(id serial, deviceid varchar(100), value int, other varchar(1000), clienttime int8, servertime int8, createdtime int8, created timestamp with time zone default(now()))")
	if err != nil {
		fmt.Println("create table error:", err)
		return
	}
	fmt.Println("Table dvdata created")
}
