package main

import (
	"fmt"

	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler/scheduler"
	"github.com/ashin-l/go-exercise/crawler/zhenai/parser"
	"github.com/ashin-l/go-exercise/crawler_distributed/config"
	itemsaver "github.com/ashin-l/go-exercise/crawler_distributed/persist/client"
	worker "github.com/ashin-l/go-exercise/crawler_distributed/worker/client"
)

func main() {
	itemChan, err := itemsaver.ItemSaver(fmt.Sprintf(":%d", config.ItemSaverPort))
	if err != nil {
		panic(err)
	}
	processor, err := worker.CreateProcessor()
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
