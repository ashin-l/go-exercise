package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ashin-l/go-exercise/logcollection/server/common"

	etcd "go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   common.AppConf.EtcdAddrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Println("create client err:", err)
	}
	defer cli.Close()

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	//resp, err := cli.Put(ctx, "name", "ashin")
	//cancle()
	//if err != nil {
	//	fmt.Println("put err:", err)
	//}
	//fmt.Println("resp:", resp)

	resp1, err := cli.Get(ctx, "name")
	cancle()
	if err != nil {
		fmt.Println("put err:", err)
	}
	fmt.Println("resp:", resp1)
}
