package fetcher

import (
	"fmt"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/crawler"
	"regexp"
	"strconv"
	"strings"
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

func ExtractSeekJobId(html *string) []string {
	regex := regexp.MustCompile(`data-job-id="(.*?)"`)
	matches := regex.FindAllStringSubmatch(*html, -1)
	res := make([]string, len(matches))
	for idx, match := range matches {
		res[idx] = match[1]
	}
	return res
}

func (fetcher *SeekFetcher) FetchContents(criteria *FetchCriteria) ([]string, error) {
	urls, err := fetcher.FetchTargetJobUrls(criteria)
	if err != nil {
		return nil, err
	}

	return fetcher.textCrawler.FetchManyPages(urls)
}

func (fetcher *SeekFetcher) FetchCriteriaToUrl(fc *FetchCriteria, pageNumber int) (string, error) {
	if len(fc.Description) == 0 || len(fc.JobTitle) == 0 || fc.Duration <= 0 || len(fc.Country) == 0 {
		return "", fmt.Errorf("one of fetch criteria is missing or invalid: expect description, job title and country be non-empty and duration be greater than 0")
	}
	if _, exist := (*fetcher.countryCode)[fc.Country]; !exist {
		return "", fmt.Errorf("country: %s not supported", fc.Country)
	}

	jobString := strings.ReplaceAll(fc.JobTitle, " ", "-")
	countryString := strings.ReplaceAll(fc.Country, " ", "-")
	return fmt.Sprintf("https://www.seek.com.%s/%s-jobs/in-%s?daterange=%d&page=%d",
		(*fetcher.countryCode)[fc.Country], jobString, countryString, fc.Duration, pageNumber), nil
}

func (*SeekFetcher) ExtractJobDescriptionUrls(html *string) []string {
	res := ExtractSeekJobId(html)
	for idx, id := range res {
		res[idx] = "https://www.seek.com.au/job/" + id
	}
	return res
}

func (*SeekFetcher) ExtractTotalJobCounts(html *string) int {
	re := regexp.MustCompile(`<span\sdata-automation="totalJobsCount">(.*?)</span>`)
	matches := re.FindStringSubmatch(*html)
	if len(matches) == 0 {
		return 0
	}
	count, _ := strconv.Atoi(matches[1])
	return count
}

func (fetcher *SeekFetcher) FetchTargetJobUrls(fc *FetchCriteria) ([]string, error) {
	url, err := fetcher.FetchCriteriaToUrl(fc, 1)
	fmt.Println(url)
	if err != nil {
		return nil, err
	}

	html, err := fetcher.htmlCrawler.FetchOnePage(url)
	if err != nil {
		return nil, err
	}

	jobFoundCount := fetcher.ExtractTotalJobCounts(&html)
	jobDescriptionUrls := fetcher.ExtractJobDescriptionUrls(&html)

	pageNumber := 2
	for len(jobDescriptionUrls) < jobFoundCount && pageNumber <= 3 { //ToDo: More adaptive way to detect stopping point instead of magic number way
		url, _ = fetcher.FetchCriteriaToUrl(fc, pageNumber)
		fmt.Println(url)
		html, err := fetcher.htmlCrawler.FetchOnePage(url)
		if err != nil {
			return nil, err
		}
		jobDescriptionUrls = append(jobDescriptionUrls, fetcher.ExtractJobDescriptionUrls(&html)...)
		pageNumber += 1
	}
	return jobDescriptionUrls, nil
}
