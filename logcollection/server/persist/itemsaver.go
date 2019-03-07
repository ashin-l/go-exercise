package persist

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ashin-l/go-exercise/logcollection/server/common"

	"github.com/olivere/elastic"
)

var escli *elastic.Client

func InitES() (err error) {
	common.Logger.Info("init elasticsearch...")
	escli, err = elastic.NewClient(elastic.SetURL("http://" + common.AppConf.ESAddr))
	return err
}

func ItemSaver(index string) (chan common.Item, error) {
	out := make(chan common.Item)

	go func() {
		itemCount := 0
		for {
			item := <-out
			fmt.Printf("Item Saver: got item #%d: %v", itemCount, item)
			itemCount++
			err := Save(index, item)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}

func Save(index string, item common.Item) error {
	if item.Type == "" {
		return errors.New("must supply Type")
	}

	indexService := escli.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err := indexService.Do(context.Background())

	if err != nil {
		return err
	}

	return nil
}
