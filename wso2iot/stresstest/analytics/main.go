package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbhost := "localhost"
	dbport := "3306"
	dbuser := "root"
	dbpassword := "p000"
	dbname := "stressdb"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	dburl := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&loc=Local"
	err := orm.RegisterDataBase("default", "mysql", dburl)
	if err != nil {
		fmt.Println(err)
	}
	//orm.DefaultTimeLoc = time.Local
	orm.RegisterModel(new(StressData))
}

type StressData struct {
	Id         int
	Owner      string
	DeviceId   string
	ClientTime int64
	Pmsensor   int
	Other      string
	ServerTime int64
	CreatedAt  string
	CreateTime int64 `orm:"-"`
}

func main() {
	var datas []StressData
	qs := orm.NewOrm().QueryTable("stress_data")
	_, err := qs.Filter("id__gt", 25000000).Limit(-1).All(&datas)
	if err != nil {
		fmt.Println(err)
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04:05.000000", datas[0].CreatedAt, time.Local)
	datas[0].CreateTime = t.UnixNano() / 1e6
	tmp := datas[0].ServerTime - datas[0].ClientTime
	smin, smax, snum := tmp, tmp, tmp
	tmp = datas[0].CreateTime - datas[0].ClientTime
	dmin, dmax, dnum := tmp, tmp, tmp
	for i := 1; i != len(datas); i++ {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05.000000", datas[i].CreatedAt, time.Local)
		datas[i].CreateTime = t.UnixNano() / 1e6
		tmp = datas[i].ServerTime - datas[i].ClientTime
		if smin > tmp {
			smin = tmp
		} else if smax < tmp {
			smax = tmp
		}
		snum += tmp
		tmp = datas[i].CreateTime - datas[i].ClientTime
		if dmin > tmp {
			dmin = tmp
		} else if dmax < tmp {
			dmax = tmp
		}
		dnum += tmp
	}
	fileName := "logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	logFile, err := os.Create(fileName)
	if err != nil {
		panic("open file error !")
	}
	// 创建一个日志对象
	infoLog := log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.SetFlags(log.Lshortfile)
	fmt.Println("数据条数：", len(datas))
	fmt.Printf("到分析服务器最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", smin, smax, int(snum)/len(datas))
	fmt.Printf("到数据库最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", dmin, dmax, int(dnum)/len(datas))
	infoLog.Println("数据条数：", len(datas))
	infoLog.Printf("到分析服务器最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", smin, smax, int(snum)/len(datas))
	infoLog.Printf("到数据库最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", dmin, dmax, int(dnum)/len(datas))
	logFile.Close()
}
