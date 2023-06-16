package main

import (
	"example.com/m/v2/internal/crawler"
	"fmt"
	"time"
)

func main() {
	duration := 30 * time.Second
	frc := crawler.NewTextCrawler(duration)
	fmt.Print(frc.FetchOnePage("https://uk.indeed.com/jobs?q=software+engineer"))
}
