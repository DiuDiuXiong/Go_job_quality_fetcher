package crawler

import (
	"example.com/m/v2/internal/crawler"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

var (
	html_cl crawler.Crawler
)

func setupHtml() {
	if html_cl == nil {
		html_cl = crawler.NewHTMLCrawler(30*time.Second, crawler.NewChromeExecutor())
	}
}

func WriteHtmlToTmp(html *string, fileName string) error {
	newPath := path.Join("../fetcher", "tmp", fileName)
	err := ioutil.WriteFile(newPath, []byte(*html), 0644)
	return err
}

func TestFetchOnePage(t *testing.T) {
	setupHtml()

	targets := []struct {
		url      string
		fileName string
		target   string
	}{

		{
			url:      "https://www.seek.com.au/SRE-jobs/in-Australia?daterange=7",
			fileName: "SeekSearch.txt",
			target:   "Seek",
		},
		{
			url:      "https://au.indeed.com/jobs?q=Software+Engineer&l=Australia&start=10&fromage=10",
			fileName: "IndeedSearch.txt",
			target:   "Indeed",
		},
	}

	for _, tt := range targets {
		html, err := html_cl.FetchOnePage(tt.url)
		if err != nil {
			t.Errorf("Fetch to %s failed, please check if crawler is up to date", tt.target)
		}

		if len(html) < 20000 {
			t.Errorf("Fetch to %s failed, likely to be detected by Cloudflare, or if response within 20000 chars, check if website worth crawling", tt.target)
		}

		if err := WriteHtmlToTmp(&html, tt.fileName); err != nil {
			t.Errorf("Failed to export %s fetch one page result to fetcher folder", tt.target)
		}

	}

}
