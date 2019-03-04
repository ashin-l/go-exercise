package main

import (
	"context"
	"fmt"
	"time"

	ETCD "go.etcd.io/etcd/clientv3"
)

func main() {
	//addrs := []string{"localhost:2379", "localhost:22379", "localhost:32379"}
	addrs := []string{"localhost:2379"}
	cli, err := ETCD.New(ETCD.Config{
		Endpoints:   addrs,
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
