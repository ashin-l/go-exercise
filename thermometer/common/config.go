package common

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/config"
)

type AppConfig struct {
	LogPath     string
	LogLevel    string
	Pubinterval int
	MsgSize     int
	DBhost      string
	DBuser      string
	DBpass      string
	DBname      string
	DBQsize     int
	DeviceNum   int
	IDprefix    string
	Tenant      string
	Nameprefix  string
}

var AppConf *AppConfig

func checkLen(val, key string) (err error) {
	if len(val) == 0 {
		err = errors.New(fmt.Sprintf("配置文件出错：%s 非法!", key))
	}
	return
}

func InitConfig(confType, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		return
	}

	AppConf = &AppConfig{}
	AppConf.LogPath = conf.String("logpath")
	if len(AppConf.LogPath) == 0 {
		AppConf.LogPath = "logs"
	}

	AppConf.LogLevel = conf.String("loglevel")
	if len(AppConf.LogLevel) == 0 {
		AppConf.LogLevel = "debug"
	}

	AppConf.DeviceNum, err = conf.Int("devicenum")
	if err != nil || AppConf.DeviceNum <= 0 {
		fmt.Println("配置文件出错：devicenum 非法!")
		return err
	}

	AppConf.Pubinterval, err = conf.Int("pubinterval")
	if err != nil {
		fmt.Println("配置文件出错：pubinterval 非法!")
		return err
	}

	AppConf.MsgSize, err = conf.Int("msgsize")
	if err != nil || AppConf.MsgSize <= 0 {
		fmt.Println("配置文件出错：msgsize 非法!")
		return err
	}

	AppConf.DBhost = conf.String("dbhost")
	err = checkLen(AppConf.DBhost, "dbhost")
	if err != nil {
		return
	}

	AppConf.DBuser = conf.String("dbuser")
	err = checkLen(AppConf.DBuser, "dbuser")
	if err != nil {
		return
	}

	AppConf.DBpass = conf.String("dbpass")

	AppConf.DBname = conf.String("dbname")
	err = checkLen(AppConf.DBname, "dbname")
	if err != nil {
		return
	}

	AppConf.DBQsize, err = conf.Int("dbqsize")
	if err != nil || AppConf.DBQsize <= 0 {
		AppConf.DBQsize = 3
	}

	AppConf.IDprefix = conf.String("idprefix")
	err = checkLen(AppConf.IDprefix, "idprefix")
	if err != nil {
		return
	}

	AppConf.Tenant = conf.String("tenant")
	err = checkLen(AppConf.IDprefix, "tenant")
	if err != nil {
		return
	}

	AppConf.Nameprefix = conf.String("nameprefix")
	err = checkLen(AppConf.IDprefix, "nameprefix")
	if err != nil {
		return
	}

	return nil
}
