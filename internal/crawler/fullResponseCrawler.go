package crawler

import "github.com/chromedp/chromedp"

type FullResponseCrawler struct {
	crawlerFunction func(string) (string, error)
}

func NewFullResponseCrawler() *FullResponseCrawler {
	actionGenerator := func(dest *string) chromedp.Action {
		return chromedp.OuterHTML("html", dest, chromedp.ByQuery)
	}
	return &FullResponseCrawler{
		crawlerFunction: generateCrawlerFunction(actionGenerator, false),
	}
}

func (frc *FullResponseCrawler) FetchOnePage(url string) (string, error) {
	return frc.crawlerFunction(url)
}
