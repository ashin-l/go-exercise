package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func createDB() {
	db_host := "localhost"
	db_user := "test11"
	db_pass := "111111"
	db_port := "5432"
	db_sslmode := "disable"
	db_name := "testdb"
	dsn := fmt.Sprintf("host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_host, db_user, db_pass, db_port, db_sslmode)
	sqlstring := fmt.Sprintf("CREATE DATABASE %s", db_name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	r, err := db.Exec(sqlstring)
	if err != nil {
		log.Println(err)
		log.Println(r)
	} else {
		log.Println("Database ", db_name, " created")
	}
	defer db.Close()

}

func createTable() {
	connStr := "postgres://test11:111111@localhost/testdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.Query("create table tb1(name varchar(20), signupdate date, age int)")

}

func insert() {
	connStr := "postgres://test11:111111@localhost/testdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	sqlstr := "insert into tb1(name, age, signupdate) values('张三', 12, '1990-1-1')"
	db.Query(sqlstr)
}

func main() {
	//createDB()
	//createTable()
	insert()
}
