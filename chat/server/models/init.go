package models

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/config"
	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool
	usermgr   *UserMgr
	clientmgr *ClientMgr
)

func init() {
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("app.conf failed!")
		os.Exit(0)
		return
	}
	dburi := conf.String("dbhost") + ":" + conf.String("dbport")
	dbpasswd := conf.String("dbpasswd")
	dbId, err := conf.Int("dbid")
	//initRedis(dburi, dbpasswd, dbId, 16, 1024, 300*time.Second)
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   1024,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", dburi, redis.DialPassword(dbpasswd), redis.DialDatabase(dbId))
		},
	}
	usermgr = NewUserMgr(pool)
	clientmgr = NewClientMgr()
	return
}

func GetConn() redis.Conn {
	return pool.Get()
}

func PutConn(conn redis.Conn) {
	conn.Close()
}
