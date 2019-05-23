package common

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/config"
)

type AppConfig struct {
	LogPath   string
	LogLevel  string
	DeviceNum int
	DBhost    string
	DBport    string
	DBuser    string
	DBpass    string
	DBname    string
}

var AppConf *AppConfig

func InitConfig(confType, filename string) error {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		return err
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
	if err != nil {
		fmt.Println("配置文件出错：devicenum 非法!")
		return err
	}

	AppConf.DBhost = conf.String("dbhost")
	if len(AppConf.DBhost) == 0 {
		err = errors.New("配置文件出错：dbhost 非法!")
		return err
	}

	AppConf.DBport = conf.String("dbport")
	if len(AppConf.DBport) == 0 {
		err = errors.New("配置文件出错：dbport 非法!")
		return err
	}

	AppConf.DBuser = conf.String("dbuser")
	if len(AppConf.DBuser) == 0 {
		err = errors.New("配置文件出错：dbuser 非法!")
		return err
	}

	AppConf.DBpass = conf.String("dbpass")

	AppConf.DBname = conf.String("dbname")
	if len(AppConf.DBname) == 0 {
		err = errors.New("配置文件出错：dbname 非法!")
		return err
	}
	return nil
}
