package fetcher

import (
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/crawler"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/fetcher"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

var (
	linkedin fetcher.Fetcher
)

func SetupLinkedinTest(chromeExe *crawler.ChromeExecutor) {
	if linkedin == nil {
		linkedin = fetcher.NewLinkedinFetcher(time.Hour, chromeExe)
	}
}

func TestLinkedinFetchCriteriaToUrl(t *testing.T) {
	SetupLinkedinTest(crawler.NewChromeExecutor())
	tests := []struct {
		name        string
		criteria    *fetcher.FetchCriteria
		startNumber int
		wantErr     bool
		wantURL     string
	}{
		{
			name: "Test_Empty_Description",
			criteria: &fetcher.FetchCriteria{
				Description: "",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "USA",
			},
			startNumber: 0,
			wantErr:     true,
		},
		{
			name: "Test_Empty_JobTitle",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "",
				Duration:    10,
				Country:     "USA",
			},
			startNumber: 0,
			wantErr:     true,
		},
		{
			name: "Test_Negative_Duration",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    -1,
				Country:     "USA",
			},
			startNumber: 0,
			wantErr:     true,
		},
		{
			name: "Test_Empty_Country",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "",
			},
			startNumber: 0,
			wantErr:     true,
		},
		{
			name: "Test_Unsupported_Country",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "Mars",
			},
			startNumber: 0,
			wantErr:     true,
		},
		{
			name: "Test_Correct_Operation",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "Australia",
			},
			startNumber: 0,
			wantErr:     false,
			wantURL:     "https://au.linkedin.com/jobs/api/seeMoreJobPostings/search?keywords=Software+Engineer&location=Australia&start=0&f_TPR=r864000",
		},
	}

	for _, tt := range tests {
		url, err := linkedin.FetchCriteriaToUrl(tt.criteria, tt.startNumber)
		if (err != nil) != tt.wantErr {
			t.Errorf("FetchCriteriaToUrl() error = %v, wantErr %v", err, tt.wantErr)
		}
		if err != nil {
			continue
		}
		if url != tt.wantURL {
			t.Errorf("FetchCriteriaToUrl() url = %s, wantURL = %s", url, tt.wantURL)
		}
	}
}

func TestLinkedinExtractIndeedJobID(t *testing.T) {
	// skipped since we directly extract job urls from search file
}

func TestLinkedinExtractJobDescriptionUrls(t *testing.T) {
	SetupLinkedinTest(crawler.NewChromeExecutor())
	filePath := path.Join("tmp", "LinkedInSearch.txt")
	html, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to open %s", filePath)
	}

	htmlStr := string(html)

	jks := linkedin.ExtractJobDescriptionUrls(&htmlStr)
	if len(jks) > 30 || len(jks) < 10 {
		t.Errorf("Extract Indeed Job ID failed, with %d results, please refer to /internal/crawler/html_test.go to see if the link it outdated, or if Job ID number changed.", len(jks))
	}
}

func TestLinkedInExtractTotalJobCounts(t *testing.T) {
	// skipped since LinkedIn use a scroll method, so if we run out of jobs, it will just return empty page, which cause empty []string that won't affect results
}

func TestLinkedinFetchTargetJobUrls(t *testing.T) {
	// just a wrapper so skip
}
