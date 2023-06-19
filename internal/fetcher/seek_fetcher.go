package fetcher

import (
	"example.com/m/v2/internal/crawler"
	"time"
)

type SeekFetcher struct {
	htmlCrawler crawler.Crawler
	textCrawler crawler.Crawler
	countryCode *map[string]string
}

func NewSeekFetcher(timeout time.Duration, chromeExecutor *crawler.ChromeExecutor) Fetcher {
	return &SeekFetcher{
		htmlCrawler: crawler.NewHTMLCrawler(timeout, chromeExecutor),
		textCrawler: crawler.NewTextCrawler(timeout, chromeExecutor),
		countryCode: &map[string]string{
			"Australia": "au",
		},
	}
}

func (*SeekFetcher) FetchContents(criteria *FetchCriteria) ([]string, error) {
	return nil, nil
}

func (*SeekFetcher) FetchCriteriaToUrl(fc *FetchCriteria, pageNumber int) (string, error) {
	return "", nil
}

func (*SeekFetcher) extractJobDescriptionUrls(html *string) []string {
	return nil
}

func (*SeekFetcher) extractTotalJobCounts(html *string) int {
	return 0
}

func (*SeekFetcher) FetchTargetJobUrls(fc *FetchCriteria) ([]string, error) {
	return nil, nil
}
