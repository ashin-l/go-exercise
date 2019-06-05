package main

import (
	"time"
	"os"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Device struct {
	Id int
	Name string
	DeviceId string
	AccessToken string
	Created time.Time
}


func createDB() {
	db_host := "192.168.152.44"
	db_user := "postgres"
	db_pass := "111"
	db_port := "5432"
	db_sslmode := "disable"
	db_name := "tbstress"
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

func connectDB() {
	connStr := "postgres://postgres:111@192.168.152.44/tbstress?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
}

func insert() {
	sqlstr := "insert into tb1(deviceid, accesstoken) values('11111', 'xxxx')"
	db.Query(sqlstr)
}

func getdata() []Device {
	sqlstr := "select * from device limit $1"
	rows, err := db.Query(sqlstr, 10)
	if err != nil {
		fmt.Println("select error:", err)
		return nil
	}
	var sdv []Device
	for rows.Next() {
		var tmp Device
		err = rows.Scan(&tmp.Id, &tmp.DeviceId, &tmp.AccessToken, &tmp.Created)
		sdv = append(sdv, tmp)
	}
	return sdv
}

func main() {
	//createDB()
	//createTable()
	connectDB()
	//insert()
	sdv := getdata()
	fmt.Println(sdv)
}
