package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"strings"

	"github.com/ashin-l/go-exercise/crawler_distributed/rpcsupport"

	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler/scheduler"
	"github.com/ashin-l/go-exercise/crawler/zhenai/parser"
	"github.com/ashin-l/go-exercise/crawler_distributed/config"
	itemsaver "github.com/ashin-l/go-exercise/crawler_distributed/persist/client"
	worker "github.com/ashin-l/go-exercise/crawler_distributed/worker/client"
)

var itemSaverHost = flag.String("itemsaver_host", "", "itemsaver host")
var workerHosts = flag.String("worker_hosts", "", "worker hosts (comma separated)")

func main() {
	flag.Parse()
	if *itemSaverHost == "" || *workerHosts == "" {
		fmt.Println("must specify ports")
		return
	}
	itemChan, err := itemsaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}
	hosts := strings.Split(*workerHosts, ",")
	pool := createClientPool(hosts)
	processor := worker.CreateProcessor(pool)
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.SimpleScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}

	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun/xian",
		Parser: engine.NewFuncParser(parser.ParseCity, config.ParseCity),
	})

	//e.Run(engine.Request{
	//	Url:        "http://www.zhenai.com/zhenghun",
	//	ParserFunc: parser.ParseCityList,
	//})
}

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err != nil {
			log.Printf("Error connecting to port %s: %v", h, err)
		} else {
			clients = append(clients, client)
			log.Printf("Connected to port %s", h)
		}
	}
	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
