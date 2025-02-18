package client

import (
	"log"

	"github.com/ashin-l/go-exercise/crawler_distributed/config"

	"github.com/ashin-l/go-exercise/crawler_distributed/rpcsupport"

	"github.com/ashin-l/go-exercise/crawler/engine"
)

func ItemSaver(host string) (chan engine.Item, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item #%d: %v", itemCount, item)
			itemCount++
			result := ""
			err = client.Call(config.ItemSaverRpc, item, &result)
			if err != nil || result != "ok" {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}
