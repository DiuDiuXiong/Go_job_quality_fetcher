package fetcher

import (
	"example.com/m/v2/internal/crawler"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type IndeedFetcher struct {
	htmlCrawler crawler.Crawler
	textCrawler crawler.Crawler
	countryCode *map[string]string
}

func NewIndeedFetcher(timeout time.Duration, chromeExecutor *crawler.ChromeExecutor) Fetcher {
	return &IndeedFetcher{
		// ToDo: potential modification here since operations from a single crawler is not parallelizable as needs to navigate to some urls and each crawler only have one tab(new context), maybe include multiple crawler
		htmlCrawler: crawler.NewHTMLCrawler(timeout, chromeExecutor),
		textCrawler: crawler.NewTextCrawler(timeout, chromeExecutor),
		countryCode: &map[string]string{
			"Australia": "au",
		},
	}
}

//ToDo: Add test cases (To examine all possible corner cases)
func (fetcher *IndeedFetcher) FetchCriteriaToUrl(fc *FetchCriteria, pageNumber int) (string, error) {
	if len(fc.Description) == 0 || len(fc.JobTitle) == 0 || fc.Duration <= 0 || len(fc.Country) == 0 {
		return "", fmt.Errorf("one of fetch criteria is missing or invalid: expect description, job title and country be non-empty and duration be greater than 0")
	}
	if _, exist := (*fetcher.countryCode)[fc.Country]; !exist {
		return "", fmt.Errorf("country: %s not supported", fc.Country)
	}
	jobString := strings.ReplaceAll(fc.JobTitle, " ", "+")
	return fmt.Sprintf("https://%s.indeed.com/jobs?q=%s&l=%s&start=%d&fromage=%d",
		(*fetcher.countryCode)[fc.Country], jobString, fc.Country, pageNumber*10, fc.Duration), nil
}

//ToDo: Add test cases
func ExtractIndeedJobJk(html *string) []string { //Need this ID for indeed query: https://au.indeed.com/viewjob?jk=<jkID>
	regex := regexp.MustCompile(`/viewjob\?jk=(.*?)&`)
	matches := regex.FindAllStringSubmatch(*html, -1)
	res := make([]string, len(matches))
	for idx, match := range matches {
		res[idx] = match[1]
	}
	return res
}

//ToDo: Add test cases
func (*IndeedFetcher) ExtractJobDescriptionUrls(html *string) []string { // urls pointing to actual job description
	res := ExtractIndeedJobJk(html)
	for idx, jk := range res {
		res[idx] = "https://au.indeed.com/viewjob?jk=" + jk
	}
	return res
}

//ToDo: Add test cases
func (*IndeedFetcher) ExtractTotalJobCounts(html *string) int {
	re := regexp.MustCompile(`<div\s+class="jobsearch-JobCountAndSortPane-jobCount[^>]+><span>(\d+) jobs</span>`)
	matches := re.FindStringSubmatch(*html)
	if len(matches) == 0 {
		return 0
	}
	count, _ := strconv.Atoi(matches[1])
	return count
}

// ToDo: Add test cases to avoid duplicates
func (fetcher *IndeedFetcher) FetchTargetJobUrls(fc *FetchCriteria) ([]string, error) {
	url, err := fetcher.FetchCriteriaToUrl(fc, 0)
	if err != nil {
		return nil, err
	}

	html, err := fetcher.htmlCrawler.FetchOnePage(url)
	if err != nil {
		return nil, err
	}

	jobFoundCount := fetcher.ExtractTotalJobCounts(&html)
	jobDescriptionUrls := fetcher.ExtractJobDescriptionUrls(&html)

	pageNumber := 1
	for len(jobDescriptionUrls) < jobFoundCount && pageNumber <= 3 { //ToDo: More adaptive way to detect stopping point instead of magic number way
		url, _ = fetcher.FetchCriteriaToUrl(fc, pageNumber)
		html, err := fetcher.htmlCrawler.FetchOnePage(url)
		if err != nil {
			return nil, err
		}
		jobDescriptionUrls = append(jobDescriptionUrls, fetcher.ExtractJobDescriptionUrls(&html)...)
		pageNumber += 1
	}
	return jobDescriptionUrls, nil
}

func (fetcher *IndeedFetcher) FetchContents(fc *FetchCriteria) ([]string, error) {

	urls, err := fetcher.FetchTargetJobUrls(fc)
	if err != nil {
		return nil, err
	}
	// ToDo: Parallel fetch here
	return fetcher.textCrawler.FetchManyPages(urls)
}
