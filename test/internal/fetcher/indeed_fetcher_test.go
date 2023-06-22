package fetcher

import (
	"example.com/m/v2/internal/crawler"
	"example.com/m/v2/internal/fetcher"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

var (
	indeed fetcher.Fetcher
)

func SetupIndeedTest(chromeExe *crawler.ChromeExecutor) {
	if indeed == nil {
		indeed = fetcher.NewIndeedFetcher(time.Hour, chromeExe)
	}
}

func TestIndeedFetchCriteriaToUrl(t *testing.T) {
	SetupIndeedTest(crawler.NewChromeExecutor())
	tests := []struct {
		name     string
		criteria *fetcher.FetchCriteria
		page     int
		wantErr  bool
		wantURL  string
	}{
		{
			name: "Test_Empty_Description",
			criteria: &fetcher.FetchCriteria{
				Description: "",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "USA",
			},
			page:    1,
			wantErr: true,
		},
		{
			name: "Test_Empty_JobTitle",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "",
				Duration:    10,
				Country:     "USA",
			},
			page:    1,
			wantErr: true,
		},
		{
			name: "Test_Negative_Duration",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    -1,
				Country:     "USA",
			},
			page:    1,
			wantErr: true,
		},
		{
			name: "Test_Empty_Country",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "",
			},
			page:    1,
			wantErr: true,
		},
		{
			name: "Test_Unsupported_Country",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "Mars",
			},
			page:    1,
			wantErr: true,
		},
		{
			name: "Test_Correct_Operation",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "Australia",
			},
			page:    1,
			wantErr: false,
			wantURL: "https://au.indeed.com/jobs?q=Software+Engineer&l=Australia&start=10&fromage=10",
		},
	}

	for _, tt := range tests {
		url, err := indeed.FetchCriteriaToUrl(tt.criteria, tt.page)
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

func TestIndeedExtractIndeedJobJk(t *testing.T) {
	filePath := path.Join("tmp", "IndeedSearch.txt")
	html, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to open %s", filePath)
	}

	htmlStr := string(html)

	jks := fetcher.ExtractIndeedJobJk(&htmlStr)
	if len(jks) > 25 || len(jks) < 10 {
		t.Errorf("Extract Indeed Job ID failed, with %d results, please refer to /internal/crawler/html_test.go to see if the link it outdated, or if Job ID number changed.", len(jks))
	}
}

func TestIndeedExtractJobDescriptionUrls(t *testing.T) {
	// will pass this since job description urls' format are tested in /internal/crawler/text_test.go
}

func TestIndeedExtractTotalJobCounts(t *testing.T) {
	SetupIndeedTest(crawler.NewChromeExecutor())
	filePath := path.Join("tmp", "IndeedSearch.txt")
	html, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to open %s", filePath)
	}

	htmlStr := string(html)
	jcs := indeed.ExtractTotalJobCounts(&htmlStr)
	if jcs == 0 {
		t.Errorf("Extract Indeed Job counts failed, check the original url in /internal/crawler/html_test.go first to see if ")
	}
}

func TestIndeedFetchTargetJobUrls(t *testing.T) {
	// just a wrapper so skip
}
