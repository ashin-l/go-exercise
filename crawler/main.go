package main

import (
	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler/scheduler"
	"github.com/ashin-l/go-exercise/crawler/zhenai/parser"
)

func main() {
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
	}

	e.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}
