package view

import (
	"os"
	"testing"

	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler/frontend/model"
	common "github.com/ashin-l/go-exercise/crawler/model"
)

func TestSearchResultView_Render(t *testing.T) {
	view := CreateSearchResultView("template.html")
	out, err := os.Create("template.test.html")
	page := model.SearchResult{}
	page.Hits = 123
	item := engine.Item{
		Url:  "http://album.zhenai.com/u/86837825",
		Type: "zhenai",
		Id:   "86837825",
		Payload: common.Profile{
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
	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}
	err = view.Render(out, page)
	if err != nil {
		panic(err)
	}
}
