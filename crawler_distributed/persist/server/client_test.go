package main

import (
	"testing"
	"time"

	"github.com/ashin-l/go-exercise/crawler_distributed/config"

	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler/model"
	"github.com/ashin-l/go-exercise/crawler_distributed/rpcsupport"
)

func TestItemSaver(t *testing.T) {
	const host = ":1234"
	go serveRpc(host, "test1")
	time.Sleep(time.Second)
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}
	item := engine.Item{
		Url:  "http://album.zhenai.com/u/86837825",
		Type: "zhenai",
		Id:   "86837825",
		Payload: model.Profile{
			Name:       "小甜",
			Gender:     "女",
			Age:        28,
			Height:     165,
			Weight:     46,
			Income:     "",
			Marriage:   "未婚",
			Education:  "大学本科",
			Occupation: "设计师",
			Hokou:      "陕西西安",
			Xinzuo:     "天蝎座",
			House:      "已购房",
			Car:        "已购车",
		},
	}
	result := ""
	err = client.Call(config.ItemSaverRpc, item, &result)
	if err != nil || result != "ok" {
		t.Errorf("result: %s; err: %s", result, err)
	}
}
