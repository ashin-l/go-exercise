package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "taosSql"
)

const (
	//CONNDB = "root:taosdata@/tcp(192.168.152.43:0)/mydb"
	CONNDB = "root:taosdata@/tcp(127.0.0.1:0)/mydb"
	LIMIT  = 100000
)

type Dvdata struct {
	Degree     float64
	Other      string
	ClientTime string
	ServerTime string
}

type Result struct {
	Min int64
	Max int64
	Sum int64
}

type DBConn struct {
	DBS        *sql.DB
	StmtSelect *sql.Stmt
}

var (
	DBQ     chan *DBConn
	DBQSIZE int
	DBConns []*DBConn
	stmt    *sql.Stmt
	num     int
	muNum   sync.Mutex
	chret   chan *Result
)

func createDBQ(selectStr string, dbconn *DBConn) {
	var err error
	dbconn.DBS, err = sql.Open("taosSql", CONNDB)
	if err != nil {
		fmt.Println("Open database error: %s\n", err)
		os.Exit(0)
	}

	dbconn.StmtSelect, err = dbconn.DBS.Prepare(selectStr)
	if err != nil {
		fmt.Println("create stmtinsert error: %s\n", err)
		os.Exit(0)
	}
}

func transtime(ts string) int64 {
	var err error
	var t time.Time
	switch len(ts) {
	case 21:
		t, err = time.Parse("2006-01-02 15:04:05.0", ts)
	case 22:
		t, err = time.Parse("2006-01-02 15:04:05.00", ts)
	default:
		t, err = time.Parse("2006-01-02 15:04:05.000", ts)
	}
	if err != nil {
		fmt.Println("format time error:", err)
		os.Exit(0)
	}
	tm := t.UnixNano() / 1e6
	return tm
}

func caltime(offset int) (subs []int64) {
	//	db, err := sql.Open("taosSql", CONNDB)
	//	if err != nil {
	//		fmt.Println("--- Open database error: %s\n", err)
	//		return
	//	}
	//
	//	stmt, err = db.Prepare("select ts, degree, other, servertime from mydb.thermometer limit ? offset ?")
	//	if err != nil {
	//		fmt.Println("create select stmt error:", err)
	//		return
	//	}

	db := <-DBQ
	rows, err := db.StmtSelect.Query(LIMIT, offset)
	if err != nil {
		fmt.Println("stmt query error:", err)
		os.Exit(0)
	}

	dv := Dvdata{}
	for rows.Next() {
		err = rows.Scan(&dv.ClientTime, &dv.Degree, &dv.Other, &dv.ServerTime)
		if err != nil {
			fmt.Println("row.Scan error:", err)
		}
		ct := transtime(dv.ClientTime)
		st := transtime(dv.ServerTime)
		tmp := st - ct
		subs = append(subs, tmp)
	}
	DBQ <- db
	return
}

func do(wg *sync.WaitGroup, offset int) {
	defer wg.Done()
	subs := caltime(offset)
	l := len(subs)
	if len(subs) == 0 {
		return
	}
	muNum.Lock()
	num += l
	muNum.Unlock()
	min, max, sum := subs[0], subs[0], subs[0]
	for i := 1; i < l; i++ {
		if subs[i] < min {
			min = subs[i]
		} else if subs[i] > max {
			max = subs[i]
		}
		sum += subs[i]
	}
	ret := &Result{}
	ret.Min, ret.Max, ret.Sum = min, max, sum
	chret <- ret
}

func main() {
	if len(os.Args) == 2 {
		size, err := strconv.Atoi(os.Args[1])
		if err != nil || size < 1 {
			fmt.Println("bad param")
			return
		}
		DBQSIZE = size
	} else {
		DBQSIZE = 1000
	}
	start := time.Now()
	db, err := sql.Open("taosSql", CONNDB)
	if err != nil {
		fmt.Println("Open database error: %s\n", err)
		return
	}
	_, err = db.Exec("use mydb")
	if err != nil {
		fmt.Println("use database error:", err)
		return
	}

	//	stmt, err = db.Prepare("select ts, degree, other, servertime from thermometer limit ? offset ?")
	//	if err != nil {
	//		fmt.Println("create select stmt error:", err)
	//		return
	//	}

	rows, err := db.Query("select count(*) from thermometer")
	if err != nil {
		fmt.Println("select count(*) error:", err)
		return
	}
	db.Close()
	DBQ = make(chan *DBConn, DBQSIZE)
	//stmt, err = db.Prepare("select ts, degree, other, servertime from mydb.thermometer limit ? offset ?")
	stStr := "select ts, degree, other, servertime from mydb.thermometer limit ? offset ?"
	for i := 0; i < DBQSIZE; i++ {
		dbconn := &DBConn{}
		createDBQ(stStr, dbconn)
		DBQ <- dbconn
		DBConns = append(DBConns, dbconn)
	}
	fmt.Println("DB init down: ", len(DBConns), len(DBQ))

	chret = make(chan *Result)
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	count = count/LIMIT + 1
	wg := &sync.WaitGroup{}

	//subs := caltime(1, 0)
	//if len(subs) == 0 {
	//	fmt.Println("no data!")
	//	return
	//}
	//min, max, sum := subs[0], subs[0], subs[0]
	for i := 0; i < count; i++ {
		offset := LIMIT * i
		wg.Add(1)
		go do(wg, offset)
	}
	go func() {
		wg.Wait()
		close(chret)
	}()
	var rets []*Result
	for v := range chret {
		rets = append(rets, v)
	}

	for _, v := range DBConns {
		v.StmtSelect.Close()
		v.DBS.Close()
	}

	l := len(rets)
	if l <= 0 {
		fmt.Println("no data!")
		return
	}
	min, max, sum := rets[0].Min, rets[0].Max, rets[0].Sum
	for i := 1; i < l; i++ {
		if rets[i].Min < min {
			min = rets[i].Min
		} else if rets[i].Max > max {
			max = rets[i].Max
		}
		sum += rets[i].Sum
	}

	total := time.Now().Sub(start)
	fmt.Println("统计用时：", total)
	fmt.Println("数据条数：", num)
	fmt.Printf("到数据库最小用时（豪秒）：%d, 最大用时：%d, 平均用时：%d\n", min, max, sum/int64(num))
	fileName := "logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	logFile, err := os.Create(fileName)
	if err != nil {
		panic("open file error !")
	}
	// 创建一个日志对象
	infoLog := log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.SetFlags(log.Lshortfile)
	infoLog.Println("数据条数：", num)
	infoLog.Printf("到数据库最小用时（豪秒）：%d, 最大用时：%d, 平均用时：%d\n", min, max, sum/int64(num))
	logFile.Close()
}
