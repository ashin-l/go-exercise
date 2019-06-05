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
	Pubinterval int
	MsgSize int
	Host string
	Username string
	Password string
	Gettoken string
	Savedevice string
	Deldevice string
	Getdevicecredentials string
	Telemetryup string
	DBhost    string
	DBport    string
	DBuser    string
	DBpass    string
	DBname    string
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


	AppConf.Host = conf.String("host")
	err = checkLen(AppConf.Host, "host")
	if err != nil {
		return
	}

	AppConf.Username = conf.String("username")
	err = checkLen(AppConf.Username, "username")
	if err != nil {
		return
	}

	AppConf.Password = conf.String("password")
	err = checkLen(AppConf.Password, "password")
	if err != nil {
		return
	}

	AppConf.Gettoken = conf.String("gettoken")
	err = checkLen(AppConf.Gettoken, "gettoken")
	if err != nil {
		return
	}
	AppConf.Gettoken = fmt.Sprintf(AppConf.Gettoken, AppConf.Host)

	AppConf.Savedevice = conf.String("savedevice")
	err = checkLen(AppConf.Savedevice, "savedevice")
	if err != nil {
		return
	}
	AppConf.Savedevice = fmt.Sprintf(AppConf.Savedevice, AppConf.Host)

	AppConf.Deldevice = conf.String("deldevice")
	err = checkLen(AppConf.Deldevice, "deldevice")
	if err != nil {
		return
	}
	AppConf.Deldevice = fmt.Sprintf(AppConf.Deldevice, AppConf.Host)

	AppConf.Getdevicecredentials = conf.String("getdevicecredentials")
	err = checkLen(AppConf.Getdevicecredentials, "getdevicecredentials")
	if err != nil {
		return
	}
	AppConf.Getdevicecredentials = fmt.Sprintf(AppConf.Getdevicecredentials, AppConf.Host, "%s")

	AppConf.Telemetryup = conf.String("telemetryup")
	err = checkLen(AppConf.Telemetryup, "getdevicecredentials")
	if err != nil {
		return
	}
	AppConf.Telemetryup = fmt.Sprintf(AppConf.Telemetryup, AppConf.Host, "%s")

	AppConf.DBhost = conf.String("dbhost")
	err = checkLen(AppConf.DBhost, "dbhost")
	if err != nil {
		return
	}

	AppConf.DBport = conf.String("dbport")
	err = checkLen(AppConf.DBport, "dbport")
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
	return
}
