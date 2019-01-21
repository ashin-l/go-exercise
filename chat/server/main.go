package main

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego/config"
)

func main() {
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("app.conf failed!")
		os.Exit(0)
		return
	}
	dburi := conf.String("dbhost") + ":" + conf.String("dbport")
	dbpasswd := conf.String("dbpasswd")
	dbId, err := conf.Int("dbid")
	initRedis(dburi, dbpasswd, dbId, 16, 1024, 300*time.Second)
	initUserMgr()
	runServer("0.0.0.0:10000")
}
