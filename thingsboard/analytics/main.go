package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

const (
	CONNDB = "postgres://%s:%s@%s/%s?sslmode=disable"
	TBNAME = "dvdata"
)

func init() {
	dbhost := "192.168.152.43"
	//dbport := "5432"
	dbuser := "postgres"
	dbpassword := "postgres"
	dbname := "tbstress"
	orm.RegisterDriver("postgres", orm.DRPostgres)
	//dburl := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&loc=Local"
	dburl := fmt.Sprintf(CONNDB, dbuser, dbpassword, dbhost, dbname)
	err := orm.RegisterDataBase("default", "postgres", dburl)
	if err != nil {
		fmt.Println(err)
	}
	//orm.DefaultTimeLoc = time.Local
	orm.RegisterModel(new(Dvdata))
}

type Dvdata struct {
	Id            int64
	DeviceID      string `orm:"column(deviceid)"`
	Value         int
	Other         string
	ClientTime    int64 `orm:"column(clienttime)"`
	ServerTime    int64 `orm:"column(servertime)"`
	Ts            string
	Transporttime int64
	Created       time.Time
	CreatedTime   int64 `orm:"column(createdtime)"`
}

func main() {
	var data Dvdata
	err := orm.NewOrm().QueryTable(TBNAME).Filter("id", 1).One(&data)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	data.CreatedTime = data.Created.UnixNano() / 1e6
	tmp := data.ServerTime - data.ClientTime
	smin, smax, snum := tmp, tmp, tmp
	tmp = data.CreatedTime - data.ClientTime
	dmin, dmax, dnum := tmp, tmp, tmp
	data.Transporttime, err = strconv.ParseInt(data.Ts, 10, 64)
	if err != nil {
		fmt.Println("ts to int error:", err)
		return
	}
	tmp = data.Transporttime - data.ClientTime
	emin, emax, enum := tmp, tmp, tmp

	qs := orm.NewOrm().QueryTable(TBNAME).Filter("id__gt", 1)
	count, _ := qs.Count()
	offset := 0
	for {
		var datas []Dvdata
		_, err := qs.Limit(100000, offset*100000).All(&datas)
		if err != nil {
			fmt.Println(err)
			break

		}
		if len(datas) == 0 {
			break
		}
		//fmt.Printf("groupid: %d, count: %d\n", offset, len(datas))
		for i := 1; i != len(datas); i++ {
			datas[i].CreatedTime = data.Created.UnixNano() / 1e6
			tmp = datas[i].ServerTime - datas[i].ClientTime
			if smin > tmp {
				smin = tmp
			} else if smax < tmp {
				smax = tmp
			}
			snum += tmp
			tmp = datas[i].CreatedTime - datas[i].ClientTime
			if dmin > tmp {
				dmin = tmp
			} else if dmax < tmp {
				dmax = tmp
			}
			dnum += tmp
			data.Transporttime, err = strconv.ParseInt(data.Ts, 10, 64)
			if err != nil {
				fmt.Println("ts to int error:", err)
				return
			}
			tmp = data.Transporttime - data.ClientTime
			if emin > tmp {
				emin = tmp
			} else if emax < tmp {
				emax = tmp
			}
			enum += tmp
		}
		offset++
	}
	count++
	fmt.Println("数据条数：", count)
	fmt.Printf("到emq服务器最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", emin, emax, enum/count)
	fmt.Printf("到分析服务器最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", smin, smax, snum/count)
	fmt.Printf("到数据库最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", dmin, dmax, dnum/count)
	fileName := "logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	logFile, err := os.Create(fileName)
	if err != nil {
		panic("open file error !")
	}
	// 创建一个日志对象
	infoLog := log.New(logFile, "[Info]", log.LstdFlags)
	infoLog.SetFlags(log.Lshortfile)
	infoLog.Println("数据条数：", count)
	infoLog.Printf("到分析服务器最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", smin, smax, snum/count)
	infoLog.Printf("到数据库最小用时（毫秒）：%d, 最大用时：%d, 平均用时：%d\n", dmin, dmax, dnum/count)
	logFile.Close()
}
