package parser

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ashin-l/go-exercise/crawler/engine"

	"github.com/ashin-l/go-exercise/crawler/model"
)

func TestParseProfile(t *testing.T) {
	contents, err := ioutil.ReadFile("profile_test_data.html")
	fmt.Println(string(contents))

	if err != nil {
		panic(err)
	}

	result := ParseProfile(contents, "http://album.zhenai.com/u/86837825", "小甜")
	actual := result.Items[0]
	fmt.Println(actual)
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

	if actual != expected {
		t.Errorf("expected item: %v; but was %v", expected, actual)
	}
}
