package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbhost := "localhost"
	dbport := "3306"
	dbuser := "root"
	dbpassword := "p000"
	dbname := "EnvMonitorDM_DB"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	dburl := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&loc=Local"
	err := orm.RegisterDataBase("default", "mysql", dburl)
	if err != nil {
		fmt.Println(err)
	}
	orm.DefaultTimeLoc = time.Local
	orm.RegisterModel(new(EnvData))
}

type EnvData struct {
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
	var datas []EnvData
	qs := orm.NewOrm().QueryTable("env_data")
	_, err := qs.All(&datas)
	if err != nil {
		fmt.Println(err)
	}
	for i := range datas {
		t, _ := time.Parse("2006-01-02 15:04:05.000000", datas[i].CreatedAt)
		datas[i].CreateTime = t.UnixNano() / 1e6
	}
	for _, v := range datas {
		fmt.Println(v.CreateTime)
	}
}
