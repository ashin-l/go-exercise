package engine

import (
	"log"

	"github.com/ashin-l/go-exercise/crawler/fetcher"
)

func Worker(r Request) (ParseResult, error) {
	//log.Printf("Fetching %s", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fecher: error "+"fetching url %s: %v", r.Url, err)
		return ParseResult{}, err
	}
	return r.Parser.Parse(body, r.Url), nil
}
