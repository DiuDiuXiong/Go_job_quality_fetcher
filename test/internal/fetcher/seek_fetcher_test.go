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
	seek fetcher.Fetcher
)

func SetupSeekTest(chromeExe *crawler.ChromeExecutor) {
	if seek == nil {
		seek = fetcher.NewSeekFetcher(30*time.Second, chromeExe)
	}
}

func TestSeekFetchCriteriaToUrl(t *testing.T) {
	SetupSeekTest(crawler.NewChromeExecutor())
	tests := []struct {
		name     string
		criteria *fetcher.FetchCriteria
		page     int
		wantErr  bool
		wantURL  string
	}{
		{
			name: "Seek_Test_Empty_Description",
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
			name: "Seek_Test_Empty_JobTitle",
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
			name: "Seek_Test_Negative_Duration",
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
			name: "Seek_Test_Empty_Country",
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
			name: "Seek_Test_Unsupported_Country",
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
			name: "Seek_Test_Correct_Operation",
			criteria: &fetcher.FetchCriteria{
				Description: "Python",
				JobTitle:    "Software Engineer",
				Duration:    10,
				Country:     "Australia",
			},
			page:    1,
			wantErr: false,
			wantURL: "https://www.seek.com.au/Software-Engineer-jobs/in-Australia?daterange=10&page=1",
		},
	}

	for _, tt := range tests {
		url, err := seek.FetchCriteriaToUrl(tt.criteria, tt.page)
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

func TestSeekExtractSeekJobId(t *testing.T) {
	filePath := path.Join("tmp", "SeekSearch.txt")
	html, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to open %s", filePath)
	}

	htmlStr := string(html)

	ids := fetcher.ExtractSeekJobId(&htmlStr)
	if len(ids) > 25 || len(ids) < 10 {
		t.Errorf("Extract Seek Job ID failed, with %d results, please refer to /internal/crawler/html_test.go to see if the link it outdated, or if Job ID number changed.", len(ids))
	}
}

func TestSeekExtractJobDescriptionUrls(t *testing.T) {
	// will pass this since job description urls' format are tested in /internal/crawler/text_test.go
}

func TestSeekExtractTotalJobCounts(t *testing.T) {
	SetupSeekTest(crawler.NewChromeExecutor())
	filePath := path.Join("tmp", "SeekSearch.txt")
	html, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to open %s", filePath)
	}

	htmlStr := string(html)
	jcs := seek.ExtractTotalJobCounts(&htmlStr)
	if jcs == 0 {
		t.Errorf("Extract Seek Job counts failed, check the original url in /internal/crawler/html_test.go first to see if ")
	}
}

func TestSeekFetchTargetJobUrls(t *testing.T) {
	// just a wrapper so skip
}
