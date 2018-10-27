package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/ashin-l/go-exercise/crawler_distributed/config"
	"github.com/ashin-l/go-exercise/crawler_distributed/rpcsupport"
	"github.com/ashin-l/go-exercise/crawler_distributed/worker"
)

func TestCrawlService(t *testing.T) {
	const host = ":9000"
	go rpcsupport.ServeRpc(host, worker.CrawlService{})
	time.Sleep(3 * time.Second)
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}

	req := worker.Request{
		Url: "http://album.zhenai.com/u/104179994",
		Parser: worker.SerializedParser{
			Name: config.ParseProfile,
			Args: "小甜",
		},
	}
	var result worker.ParseResult
	err = client.Call(config.CrawlServiceRpc, req, &result)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}
