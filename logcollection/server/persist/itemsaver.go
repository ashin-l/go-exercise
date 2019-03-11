package persist

import (
	"context"
	"fmt"
	"log"

	"github.com/ashin-l/go-exercise/logcollection/server/common"

	"github.com/olivere/elastic"
)

var escli *elastic.Client

func InitES() (err error) {
	common.Logger.Info("init elasticsearch...")
	addr := "http://" + common.AppConf.ESAddr
	fmt.Println(addr)
	escli, err = elastic.NewClient(elastic.SetURL(addr), elastic.SetSniff(false))
	return err
}

func ItemSaver(index string) chan common.Item {
	out := make(chan common.Item)

	go func() {
		itemCount := 0
		for {
			item := <-out
			fmt.Printf("Item Saver: got item #%d: %v", itemCount, item)
			itemCount++
			err := save(index, item)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()
	return out
}

func save(index string, item common.Item) error {
	indexService := escli.Index().
		Index(index).
		Type(index).
		BodyJson(item)
	_, err := indexService.Do(context.Background())

	if err != nil {
		return err
	}

	return nil
}
