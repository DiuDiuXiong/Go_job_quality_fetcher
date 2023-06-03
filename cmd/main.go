package main

import (
	"example.com/m/v2/internal/crawler"
	"fmt"
)

func main() {
	frc := crawler.NewFullResponseCrawler()
	fmt.Print(frc.FetchOnePage("https://uk.indeed.com/jobs?q=software+engineer"))
}
