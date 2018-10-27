package main

import (
	"fmt"
	"log"

	"github.com/ashin-l/go-exercise/crawler_distributed/config"
	"github.com/ashin-l/go-exercise/crawler_distributed/rpcsupport"
	"github.com/ashin-l/go-exercise/crawler_distributed/worker"
)

func main() {
	log.Fatal(rpcsupport.ServeRpc(fmt.Sprintf(":%d", config.WorkerPort0), worker.CrawlService{}))
}
