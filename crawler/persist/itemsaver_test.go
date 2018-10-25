package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ashin-l/go-exercise/crawler/engine"

	"github.com/olivere/elastic"

	"github.com/ashin-l/go-exercise/crawler/model"
)

func TestSave(t *testing.T) {
	expected := engine.Item{
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
	// TODO: Try to start up elastic search
	// here using docker go client
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	const index = "dating_test"
	err = Save(client, index, expected)
	if err != nil {
		panic(err)
	}

	resp, err := client.Get().Index(index).Type(expected.Type).Id(expected.Id).Do(context.Background())

	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", *resp.Source)
	var actual engine.Item
	json.Unmarshal(*resp.Source, &actual)
	fmt.Println(actual)
	if err != nil {
		panic(err)
	}

	actualProfile, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile

	if actual != expected {
		t.Errorf("got %v; expected %v", actual, expected)
	}
}
