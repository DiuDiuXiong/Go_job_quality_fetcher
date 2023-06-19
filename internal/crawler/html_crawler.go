package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

type HTMLCrawler struct {
	timeout    time.Duration
	chromeExec *ChromeExecutor
}

func NewHTMLCrawler(timeout time.Duration, chromeExecutor *ChromeExecutor) *HTMLCrawler {
	return &HTMLCrawler{timeout: timeout, chromeExec: chromeExecutor}
}

func (cl *HTMLCrawler) FetchOnePage(url string) (string, error) {
	ctx, cancel := cl.chromeExec.NewContext()
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, cl.timeout)
	defer cancel()
	var res string
	var action chromedp.Action = chromedp.OuterHTML("html", &res, chromedp.ByQuery)
	if err := crawWithAction(ctx, url, &action); err != nil {
		return "", err
	}
	return res, nil

}

func (cl *HTMLCrawler) FetchManyPages(urls []string) ([]string, error) {
	ctx, cancel := cl.chromeExec.NewContext()
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Duration(int(cl.timeout)*len(urls)))
	defer cancel()
	res := make([]string, len(urls))

	for idx, url := range urls {
		var action chromedp.Action = chromedp.OuterHTML("html", &res[idx], chromedp.ByQuery)
		if err := crawWithAction(ctx, url, &action); err != nil {
			return nil, err
		}
	}

	return res, nil
}
