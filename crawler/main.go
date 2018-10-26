package main

import (
	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler/persist"
	"github.com/ashin-l/go-exercise/crawler/scheduler"
	"github.com/ashin-l/go-exercise/crawler/zhenai/parser"
)

func main() {
	itemChan, err := persist.ItemSaver("dating_profile")
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 100,
		ItemChan:    itemChan,
	}

	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun/xian",
		Parser: engine.NewFuncParser(parser.ParseCity, "ParseCity"),
	})

	//e.Run(engine.Request{
	//	Url:        "http://www.zhenai.com/zhenghun",
	//	Parser: engine.NewFuncParser(parser.ParseCityList, "ParseCityList"),
	//})
}
