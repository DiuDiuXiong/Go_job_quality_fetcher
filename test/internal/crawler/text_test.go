package crawler

import (
	"example.com/m/v2/internal/crawler"
	"strings"
	"testing"
	"time"
)

var (
	text_cl crawler.Crawler
)

func setupText() {
	if text_cl == nil {
		text_cl = crawler.NewTextCrawler(30*time.Second, crawler.NewChromeExecutor())
	}
}

func TestFetchJobDescriptionURLFormat(t *testing.T) {
	setupText()
	tests := []struct {
		target          string
		urls            string
		notFoundMessage []string
	}{
		{
			target:          "Seek",
			urls:            "https://www.seek.com.au/job/68309011",
			notFoundMessage: []string{"We couldn’t find that page"},
		},
		{
			target:          "Indeed",
			urls:            "https://au.indeed.com/viewjob?jk=87d5a9a8a5cd6a3f",
			notFoundMessage: []string{"The page you requested could not be found.", "We can’t find this page"},
		},
	}

	for _, tt := range tests {
		pageText, err := text_cl.FetchOnePage(tt.urls)
		if err != nil {
			t.Errorf("Failed to fetch job description page of %s, here is the target url: %s", tt.target, tt.urls)
		}

		for _, NotFoundStr := range tt.notFoundMessage {
			if strings.Contains(pageText, NotFoundStr) {
				t.Errorf("Seems the fetch url format is broken, please check if job deprecated or job link format deprecated")
			}
		}
	}
}
