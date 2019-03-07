package common

import (
	"errors"

	"github.com/astaxie/beego/config"
)

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

	AppConf.KafkaAddrs = conf.Strings("kafka_addrs")
	if len(AppConf.KafkaAddrs) == 0 {
		err = errors.New("must config kafka_addrs!")
		return err
	}

	AppConf.ESAddr = conf.String("elastic_addr")
	if len(AppConf.ESAddr) == 0 {
		err = errors.New("must config elastic_addr!")
		return err
	}

	AppConf.EtcdAddrs = conf.Strings("etcd_addrs")
	if len(AppConf.EtcdAddrs) == 0 {
		err = errors.New("must config etcd_addrs!")
		return err
	}
	return nil
}
