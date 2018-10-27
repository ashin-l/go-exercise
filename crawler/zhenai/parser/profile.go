package parser

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ashin-l/go-exercise/crawler/engine"
	"github.com/ashin-l/go-exercise/crawler/model"
	"github.com/ashin-l/go-exercise/crawler_distributed/config"
)

var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var heightRe = regexp.MustCompile(`<td><span class="label">身高：</span>([\d]+)CM</td>`)
var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">([\d]+)KG</span></td>`)
var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([0-9\-]元)</td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var occupationRe = regexp.MustCompile(`<td><span class="label">职业： </span>([^<]+)</td>`)
var hokouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var xinzuoRe = regexp.MustCompile(`<td><span class="label">星座：</span>([^<]+)</td>`)
var houseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var carRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
var idUrlRe = regexp.MustCompile(`http://album.zhenai.com/u/([\d]+)`)

func parseProfile(contents []byte, url string, name string) engine.ParseResult {
	profile := model.Profile{}
	profile.Name = name
	profile.Gender = extractString(contents, genderRe)
	age, err := strconv.Atoi(extractString(contents, ageRe))
	if err != nil {
		log.Print("Error: ", err)
	}
	profile.Age = age
	height, err := strconv.Atoi(extractString(contents, heightRe))
	if err != nil {
		log.Print("Error: ", err)
	}
	profile.Height = height
	weight, err := strconv.Atoi(extractString(contents, weightRe))
	if err != nil {
		log.Print("Error: ", err)
	}
	profile.Weight = weight
	profile.Income = extractString(contents, incomeRe)
	profile.Marriage = extractString(contents, marriageRe)
	profile.Education = extractString(contents, educationRe)
	profile.Occupation = extractString(contents, occupationRe)
	profile.Hokou = extractString(contents, hokouRe)
	profile.Xinzuo = extractString(contents, xinzuoRe)
	profile.House = extractString(contents, houseRe)
	profile.Car = extractString(contents, carRe)
	result := engine.ParseResult{
		Items: []engine.Item{{
			Url:     url,
			Type:    "zhenai",
			Id:      extractString([]byte(url), idUrlRe),
			Payload: profile,
		},
		},
	}
	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		fmt.Println(match[1])
		return string(match[1])
	} else {
		return ""
	}
}

type ProfileParser struct {
	userName string
}

func (p *ProfileParser) Parse(contents []byte, url string) engine.ParseResult {
	return parseProfile(contents, url, p.userName)
}

func (p *ProfileParser) Serialize() (name string, args interface{}) {
	return config.ParseProfile, p.userName
}

func NewProfileParser(name string) *ProfileParser {
	return &ProfileParser{userName: name}
}
