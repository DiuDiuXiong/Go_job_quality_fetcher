package fetcher

import (
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/crawler"
	"time"
)

type LinkedinFetcher struct {
	htmlCrawler crawler.Crawler
	textCrawler crawler.Crawler
	countryCode *map[string]string
}

func NewLinkedinFetcher(timeout time.Duration, chromeExecutor *crawler.ChromeExecutor) Fetcher {
	return &SeekFetcher{
		htmlCrawler: crawler.NewHTMLCrawler(timeout, chromeExecutor),
		textCrawler: crawler.NewTextCrawler(timeout, chromeExecutor),
		countryCode: &map[string]string{
			"Australia": "au",
		},
	}
}

func (*LinkedinFetcher) FetchContents(criteria *FetchCriteria) ([]string, error) {
	return nil, nil
}

func (*LinkedinFetcher) FetchCriteriaToUrl(fc *FetchCriteria, pageNumber int) (string, error) {
	return "", nil
}

func (*LinkedinFetcher) ExtractJobDescriptionUrls(html *string) []string {
	return nil
}
func (*LinkedinFetcher) ExtractTotalJobCounts(html *string) int {
	return 0
}
func (*LinkedinFetcher) FetchTargetJobUrls(fc *FetchCriteria) ([]string, error) {
	return nil, nil
}
