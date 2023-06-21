package fetcher

import (
	"example.com/m/v2/internal/crawler"
	"example.com/m/v2/internal/fetcher"
	"testing"
	"time"
)

var (
	f fetcher.Fetcher
)

func SetupIndeedTest(chromeExe *crawler.ChromeExecutor) {
	if f == nil {
		f = fetcher.NewIndeedFetcher(time.Hour, chromeExe)
	}
}

func TestFetchCriteriaToUrl(t *testing.T) {
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
		url, err := f.FetchCriteriaToUrl(tt.criteria, tt.page)
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

func TestExtractIndeedJobJk(t *testing.T) {
}

func TestExtractJobDescriptionUrls(t *testing.T) {

}

func TestExtractTotalJobCounts(t *testing.T) {

}

func TestFetchTargetJobUrls(t *testing.T) {

}
