package main

import (
	"example.com/m/v2/internal/crawler"
	"fmt"
)

type windowSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func main() {
	frc := crawler.NewFullResponseCrawler()
	fmt.Print(frc.FetchOnePage("https://uk.indeed.com/jobs?q=software+engineer"))
}
