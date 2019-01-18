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
	var data StressData
	orm.NewOrm().QueryTable("stress_data").Filter("id", 62500001).One(&data)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05.000000", data.CreatedAt, time.Local)
	data.CreateTime = t.UnixNano() / 1e6
	tmp := data.ServerTime - data.ClientTime
	smin, smax, snum := tmp, tmp, tmp
	tmp = data.CreateTime - data.ClientTime
	dmin, dmax, dnum := tmp, tmp, tmp

	qs := orm.NewOrm().QueryTable("stress_data").Filter("id__gt", 62500001)
	count, _ := qs.Count()
	offset := 0
	for {
		var datas []StressData
		_, err := qs.Limit(10000000, offset*10000000).All(&datas)
		if err != nil {
			fmt.Println(err)
			break

		}
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
		offset++
	}
	fileName := "logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	logFile, err := os.Create(fileName)
	if err != nil {
		panic("open file error !")
	}
	// 创建一个日志对象
	infoLog := log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.SetFlags(log.Lshortfile)
	fmt.Println("数据条数：", count)
	fmt.Printf("到分析服务器最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", smin, smax, snum/count)
	fmt.Printf("到数据库最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", dmin, dmax, dnum/count)
	infoLog.Println("数据条数：", count)
	infoLog.Printf("到分析服务器最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", smin, smax, snum/count)
	infoLog.Printf("到数据库最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", dmin, dmax, dnum/count)
	logFile.Close()
}
