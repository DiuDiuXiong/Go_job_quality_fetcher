package fetcher

import (
	"fmt"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/crawler"
	"net/url"
	"regexp"
	"time"
)

type LinkedinFetcher struct {
	htmlCrawler crawler.Crawler
	textCrawler crawler.Crawler
	countryCode *map[string]string
}

func NewLinkedinFetcher(timeout time.Duration, chromeExecutor *crawler.ChromeExecutor) Fetcher {
	return &LinkedinFetcher{
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

func (fetcher *LinkedinFetcher) FetchCriteriaToUrl(fc *FetchCriteria, startNumber int) (string, error) { // LinkedIn use a scroll approach for more jobs, each time load 25 jobs
	if len(fc.Description) == 0 || len(fc.JobTitle) == 0 || fc.Duration <= 0 || len(fc.Country) == 0 {
		return "", fmt.Errorf("one of fetch criteria is missing or invalid: expect description, job title and country be non-empty and duration be greater than 0")
	}

	if _, exist := (*fetcher.countryCode)[fc.Country]; !exist {
		return "", fmt.Errorf("country: %s not supported", fc.Country)
	}

	dayToSec := func(d int) int { return d * 86400 }

	jobTitleEncoded := url.QueryEscape(fc.JobTitle)
	countryEncoded := url.QueryEscape(fc.Country)

	return fmt.Sprintf("https://%s.linkedin.com/jobs/api/seeMoreJobPostings/search?keywords=%s&location=%s&start=%d&f_TPR=r%d",
		(*fetcher.countryCode)[fc.Country], jobTitleEncoded, countryEncoded, startNumber, dayToSec(fc.Duration)), nil
}

func (*LinkedinFetcher) ExtractJobDescriptionUrls(html *string) []string {
	re := regexp.MustCompile(`<a class="base-card__full-link absolute top-0 right-0 bottom-0 left-0 p-0 z-\[2\]" href="(.*?)"`)
	matches := re.FindAllStringSubmatch(*html, -1)

	var results []string
	for _, match := range matches {
		if len(match) >= 2 {
			results = append(results, match[1])
		}
	}

	return results
}

func (*LinkedinFetcher) ExtractTotalJobCounts(html *string) int { // No much value here since we query the api, when we reach a startNumber exceed maximum jobs, it will just return empty page so we know we reach limit
	return 0
}

func (*LinkedinFetcher) FetchTargetJobUrls(fc *FetchCriteria) ([]string, error) {
	return nil, nil
}
