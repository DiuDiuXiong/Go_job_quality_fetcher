package main

import (
	"example.com/m/v2/internal/crawler"
	"example.com/m/v2/internal/fetcher"
	"fmt"
	"time"
)

func main() {
	duration := 30 * time.Second
	chromeExecutor := crawler.NewChromeExecutor()
	/*
		indeed := fetcher.NewIndeedFetcher(duration, chromeExecutor)

		jobDescriptions, _ := indeed.FetchContents(&fetcher.FetchCriteria{
			JobTitle:    "SRE",
			Duration:    50,
			Description: "-",
			Country:     "Australia",
		})

		for idx, jobDescription := range jobDescriptions {
			filename := fmt.Sprintf("%d.txt", idx) // Create the filename using idx

			err := ioutil.WriteFile(filename, []byte(jobDescription), 0644)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Created %s and wrote job description to it.\n", filename)
		}
	*/
	//htmlF := crawler.NewHTMLCrawler(duration, chromeExecutor)

	//fmt.Println(htmlF.FetchOnePage("https://www.seek.com.au/Tax-jobs"))
	seek := fetcher.NewSeekFetcher(duration, chromeExecutor)
	jobContents, _ := seek.FetchContents(&fetcher.FetchCriteria{
		JobTitle:    "SRE",
		Duration:    50,
		Description: "-",
		Country:     "Australia",
	})
	fmt.Println(jobContents)

}
