package client

import (
	"net/rpc"

	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler_distributed/config"
	"github.com/ashin-l/go-exercise/crawler_distributed/worker"
)

func CreateProcessor(clientChan chan *rpc.Client) engine.Processor {
	return func(req engine.Request) (engine.ParseResult, error) {
		sReq := worker.SerializeRequest(req)
		var sResult worker.ParseResult
		client := <-clientChan
		err := client.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializeParseResult(sResult), nil
	}
}
