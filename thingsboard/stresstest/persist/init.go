package persist

import (
	"database/sql"
	"fmt"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/common"
	_ "github.com/lib/pq"
)

var tbdb *sql.DB
var st *sql.Stmt
var delst *sql.Stmt
const CONNDB = "postgres://%s:%s@%s/%s?sslmode=disable"

func InitDB() (err error) {
	connStr := fmt.Sprintf(CONNDB, common.AppConf.DBuser, common.AppConf.DBpass, common.AppConf.DBhost, common.AppConf.DBname)
	tbdb, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("connect database error:", err)
		return
	}
	sqlstr := "insert into device(id, name, deviceid, accesstoken) values($1, $2, $3, $4)"
	st, err = tbdb.Prepare(sqlstr)
	if err != nil {
		fmt.Println("create insert stmt error:", err)
	}
	sqlstr = "delete from device where id = $1"
	delst, err = tbdb.Prepare(sqlstr)
	if err != nil {
		fmt.Println("create del stmt error:", err)
	}
	return
}
func Syncdb() {
	createDB()
	createTable()
}

func createDB() {
	dsn := fmt.Sprintf("host=%s  user=%s  password=%s  port=%s  sslmode=disable", common.AppConf.DBhost, common.AppConf.DBuser, common.AppConf.DBpass, common.AppConf.DBport)
	sqlstr := fmt.Sprintf("CREATE DATABASE %s", common.AppConf.DBname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("sql open error:", err)
		return
	}
	defer db.Close()
	_, err = db.Exec(sqlstr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Database ", common.AppConf.DBname, " created")
}

func createTable() {
	connStr := fmt.Sprintf(CONNDB, common.AppConf.DBuser, common.AppConf.DBpass, common.AppConf.DBhost, common.AppConf.DBname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("database open error:", err)
		return
	}
	defer db.Close()
	_, err = db.Query("create table device(id int, name varchar(100), deviceid varchar(100), accesstoken varchar(100), created date default(now()), primary key(id))")
	if err != nil {
		fmt.Println("create table error:", err)
		return
	}
	fmt.Println("Table device created")
}
