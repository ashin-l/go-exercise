package main

import (
	"fmt"

	"github.com/ashin-l/go-exercise/logcollection/server/scheduler"

	"github.com/ashin-l/go-exercise/logcollection/server/persist"

	"github.com/ashin-l/go-exercise/logcollection/server/common"

	"github.com/ashin-l/go-exercise/logcollection/server/engine"
	"github.com/ashin-l/go-exercise/logcollection/server/fetcher"
)

func main() {
	err := common.InitConfig("ini", "app.conf")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(common.AppConf)

	err = common.InitLogger()
	if err != nil {
		panic(err)
		return
	}

	err = fetcher.InitKafka()
	if err != nil {
		fmt.Printf("init kafka failed, error: %v", err)
		return
	}

	err = persist.InitES()
	if err != nil {
		fmt.Printf("init elasticSearch failed, error: %v", err)
		return
	}

	fmt.Println("Collection server start...")
	topics := []string{"testlog"}
	e := engine.ConcurrentEngine{
		Topics:      topics,
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 30,
	}

	e.Run()
	signal := make(chan struct{})
	<-signal
}
