package crawler

import "github.com/chromedp/chromedp"

type TextOnlyCrawler struct {
	crawlerFunction func(string) (string, error)
}

func NewTextOnlyCrawler() *TextOnlyCrawler {
	actionGenerator := func(dest *string) chromedp.Action {
		return chromedp.Text("html", dest, chromedp.ByQuery)
	}

	return &TextOnlyCrawler{crawlerFunction: generateCrawlerFunction(actionGenerator)}
}

func (toc *TextOnlyCrawler) FetchOnePage(url string) (string, error) {
	return toc.crawlerFunction(url)
}
