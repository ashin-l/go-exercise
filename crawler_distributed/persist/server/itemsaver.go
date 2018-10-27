package main

import (
	"fmt"

	"github.com/ashin-l/go-exercise/crawler_distributed/config"
	"github.com/ashin-l/go-exercise/crawler_distributed/persist"
	"github.com/ashin-l/go-exercise/crawler_distributed/rpcsupport"
	"github.com/olivere/elastic"
)

func main() {
	err := serveRpc(fmt.Sprintf(":%d", config.ItemSaverPort), config.ElasticIndex)
	if err != nil {
		panic(err)
	}
}

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
