package persist

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	_ "taosSql"
	"time"

	"github.com/ashin-l/go-exercise/thermometer/common"
)

const (
	CONNDB = "%s:%s@/tcp(%s)/%s"
	DRIVER = "taosSql"
	STable = "thermometer"
)

var (
	DBConns       []*common.DBConn
	DBQ           chan *common.DBConn
	GlobalDB      *sql.DB
	stCreateTable *sql.Stmt
)

func InitDB() {
	var err error
	connStr := fmt.Sprintf(CONNDB, common.AppConf.DBuser, common.AppConf.DBpass, common.AppConf.DBhost, common.AppConf.DBname)
	GlobalDB, err = sql.Open(DRIVER, connStr)
	if err != nil {
		fmt.Println("Open database error: %s\n", err)
		os.Exit(0)
	}
	_, err = GlobalDB.Exec("use " + common.AppConf.DBname)
	checkErr(err)

	stStr := "create table ? using " + STable + " tags ('?', '?')"
	stCreateTable, err = GlobalDB.Prepare(stStr)
	if err != nil {
		fmt.Println("create stmtinsert error: %s\n", err)
		os.Exit(0)
	}

	DBQ = make(chan *common.DBConn, common.AppConf.DBQsize)
	stStr = "insert into ? values(?, ?, '?', ?)"
	for i := 0; i < common.AppConf.DBQsize; i++ {
		dbconn := &common.DBConn{}
		createDBQ(connStr, stStr, dbconn)
		DBQ <- dbconn
		DBConns = append(DBConns, dbconn)
	}
	fmt.Println("DB init down: ", len(DBConns), len(DBQ))
}

func createDBQ(connStr string, insertStr string, dbconn *common.DBConn) {
	var err error
	dbconn.DBS, err = sql.Open(DRIVER, connStr)
	if err != nil {
		fmt.Println("Open database error: %s\n", err)
		os.Exit(0)
	}

	_, err = dbconn.DBS.Exec("use " + common.AppConf.DBname)
	checkErr(err)

	dbconn.StmtInsert, err = dbconn.DBS.Prepare(insertStr)
	if err != nil {
		fmt.Println("create stmtinsert error: %s\n", err)
		os.Exit(0)
	}
}

func CloseDB() {
	stCreateTable.Close()
	GlobalDB.Close()
	for _, v := range DBConns {
		v.StmtInsert.Close()
		v.DBS.Close()
	}
}

func Syncdb() {
	// open connect to taos server
	connStr := fmt.Sprintf(CONNDB, common.AppConf.DBuser, common.AppConf.DBpass, common.AppConf.DBhost, common.AppConf.DBname)
	var err error
	GlobalDB, err = sql.Open(DRIVER, connStr)
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
		return
	}
	defer GlobalDB.Close()

	dropDatabase()
	createDatabase()
	useDatabase()
	createStable()

	//insert_data(db, DBname)
	//select_data(db, DBname)

	//dbnameStmt := "dbnameStmt"
	//stableStmt := "stableStmt"
	//drop_database_stmt(db, dbnameStmt)
	//create_database_stmt(db, dbnameStmt)
	//use_database_stmt(db, dbnameStmt)
	//create_table_stmt(db, stableStmt)
	//insert_data_stmt()
	//select_data_stmt(db, stableStmt)

	fmt.Printf("\n======== end demo test ========\n")
}

func GetDevices() (dvs []common.Device, err error) {
	var rows *sql.Rows
	rows, err = GlobalDB.Query("select tbname, tenant, name from ? where tenant='?'", STable, common.AppConf.Tenant)
	if err != nil {
		fmt.Println("query error")
		return
	}

	for rows.Next() {
		var dv common.Device
		var id []byte

		err = rows.Scan(&id, &dv.TenantID, &dv.Name)
		id = bytes.Trim(id, "\x00")
		dv.DeviceID = string(id)
		checkErr(err)
		dvs = append(dvs, dv)
	}
	return
}

func dropDatabase() {
	st := time.Now().Nanosecond()
	res, err := GlobalDB.Exec("drop database " + common.AppConf.DBname)
	checkErr(err)

	affectd, err := res.RowsAffected()
	checkErr(err)

	et := time.Now().Nanosecond()

	fmt.Printf("drop database result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func createDatabase() {
	st := time.Now().Nanosecond()
	// create database
	res, err := GlobalDB.Exec("create database " + common.AppConf.DBname)
	checkErr(err)

	affectd, err := res.RowsAffected()
	checkErr(err)

	et := time.Now().Nanosecond()

	fmt.Printf("create database result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)

	return
}

func useDatabase() {
	GlobalDB.Exec("use " + common.AppConf.DBname)
}

func createStable() {
	st := time.Now().Nanosecond()
	// create table
	res, err := GlobalDB.Exec("create table " + STable + " (ts timestamp, degree double, other binary(40), servertime timestamp) tags(tenant binary(15), name binary(15))")
	checkErr(err)

	affectd, err := res.RowsAffected()
	checkErr(err)

	et := time.Now().Nanosecond()
	fmt.Printf("create table result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func insert_data(db *sql.DB, stable string) {
	st := time.Now().Nanosecond()
	// insert data
	res, err := db.Exec("insert into " + stable +
		" values (now, 100, 'beijing', 10,  true, 'one', 123.456, 123.456)" +
		" (now+1s, 101, 'shanghai', 11, true, 'two', 789.123, 789.123)" +
		" (now+2s, 102, 'shenzhen', 12,  false, 'three', 456.789, 456.789)")

	checkErr(err)

	affectd, err := res.RowsAffected()
	checkErr(err)

	et := time.Now().Nanosecond()
	fmt.Printf("insert data result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func select_data(db *sql.DB, stable string) {
	st := time.Now().Nanosecond()

	rows, err := db.Query("select * from ? ", stable) // go text mode
	checkErr(err)

	fmt.Printf("%10s%s%8s %5s %9s%s %s %8s%s %7s%s %8s%s %4s%s %5s%s\n", " ", "ts", " ", "id", " ", "name", " ", "len", " ", "flag", " ", "notes", " ", "fv", " ", " ", "dv")
	var affectd int
	for rows.Next() {
		var ts string
		var name string
		var id int
		var len int8
		var flag bool
		var notes string
		var fv float32
		var dv float64

		err = rows.Scan(&ts, &id, &name, &len, &flag, &notes, &fv, &dv)
		checkErr(err)

		fmt.Printf("%s\t", ts)
		fmt.Printf("%d\t", id)
		fmt.Printf("%10s\t", name)
		fmt.Printf("%d\t", len)
		fmt.Printf("%t\t", flag)
		fmt.Printf("%s\t", notes)
		fmt.Printf("%06.3f\t", fv)
		fmt.Printf("%09.6f\n", dv)

		affectd++
	}

	et := time.Now().Nanosecond()
	fmt.Printf("insert data result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func drop_database_stmt(db *sql.DB, dbname string) {
	st := time.Now().Nanosecond()
	// drop test db
	stmt, err := db.Prepare("drop database ?")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(dbname)
	checkErr(err)

	affectd, err := res.RowsAffected()
	checkErr(err)

	et := time.Now().Nanosecond()
	fmt.Printf("drop database result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func create_database_stmt(db *sql.DB, dbname string) {
	st := time.Now().Nanosecond()
	// create database
	//var stmt interface{}
	stmt, err := db.Prepare("create database ?")
	checkErr(err)

	//var res driver.Result
	res, err := stmt.Exec(dbname)
	checkErr(err)

	//fmt.Printf("Query OK, %d row(s) affected()", res.RowsAffected())
	affectd, err := res.RowsAffected()
	checkErr(err)

	et := time.Now().Nanosecond()
	fmt.Printf("create database result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func use_database_stmt(db *sql.DB, dbname string) {
	st := time.Now().Nanosecond()
	// create database
	//var stmt interface{}
	stmt, err := db.Prepare("use " + dbname)
	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	affectd, err := res.RowsAffected()
	checkErr(err)

	et := time.Now().Nanosecond()
	fmt.Printf("use database result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func CreateTable(id, tenant, name string) error {
	_, err := stCreateTable.Exec(id, tenant, name)
	if err != nil {
		return err
	}

	return nil
}

func InsertData(tbname string, degree float32, other string) (err error) {
	select {
	case <-time.After(time.Second):
		fmt.Println("insert timeout!")
		err = errors.New("insert timeout")
	case db := <-DBQ:
		st := time.Now().UnixNano() / 1e6
		_, err = db.StmtInsert.Exec(tbname, st, degree, other, "now")
		DBQ <- db
		l := len(DBQ)
		if l < 3 {
			fmt.Println("low dbq ", l)
			common.Logger.Info("low dbq %d", l)
		}
	}
	return
}

func select_data_stmt(db *sql.DB, stable string) {
	st := time.Now().Nanosecond()

	stmt, err := db.Prepare("select ?, ?, ?, ?, ?, ?, ?, ? from ?") // go binary mode
	checkErr(err)

	rows, err := stmt.Query("ts", "id", "name", "len", "flag", "notes", "fv", "dv", stable)
	checkErr(err)

	fmt.Printf("%10s%s%8s %5s %9s%s %s %8s%s %7s%s %8s%s %11s%s %14s%s\n", " ", "ts", " ", "id", " ", "name", " ", "len", " ", "flag", " ", "notes", " ", "fv", " ", " ", "dv")
	var affectd int
	for rows.Next() {
		var ts string
		var name string
		var id int
		var len int8
		var flag bool
		var notes string
		var fv float32
		var dv float64

		err = rows.Scan(&ts, &id, &name, &len, &flag, &notes, &fv, &dv)
		//fmt.Println("start scan fields from row.rs, &fv:", &fv)
		//err = rows.Scan(&fv)
		checkErr(err)

		fmt.Printf("%s\t", ts)
		fmt.Printf("%d\t", id)
		fmt.Printf("%10s\t", name)
		fmt.Printf("%d\t", len)
		fmt.Printf("%t\t", flag)
		fmt.Printf("%s\t", notes)
		fmt.Printf("%06.3f\t", fv)
		fmt.Printf("%09.6f\n", dv)

		affectd++

	}

	et := time.Now().Nanosecond()
	fmt.Printf("insert data result:\n %d row(s) affectd (%6.6fs)\n\n", affectd, (float32(et-st))/1E9)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
