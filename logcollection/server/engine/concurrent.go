package engine

import (
	"fmt"

	"github.com/ashin-l/go-exercise/logcollection/server/common"
	"github.com/ashin-l/go-exercise/logcollection/server/fetcher"
	"github.com/ashin-l/go-exercise/logcollection/server/persist"

	"github.com/ashin-l/go-exercise/logcollection/server/scheduler"
)

type ConcurrentEngine struct {
	Topics      []string
	Scheduler   *scheduler.QueuedScheduler
	WorkerCount int
}

func (e *ConcurrentEngine) Run() {
	for _, topic := range e.Topics {
		msgChan, err := fetcher.Fetch(topic)
		if err != nil {
			common.Logger.Error("error topic %s, %v", topic, err)
		}
		itemChan := persist.ItemSaver(topic)
		go e.work(msgChan, itemChan)
	}
}

func (e *ConcurrentEngine) work(msgChan chan string, itemChan chan common.Item) {
	for {
		fmt.Println("work...")
		msg := <-msgChan
		item := common.Item{Content: msg}
		fmt.Println(item)
		itemChan <- item
	}
}
