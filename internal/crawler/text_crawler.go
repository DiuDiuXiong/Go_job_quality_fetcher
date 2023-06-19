package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

type TextCrawler struct {
	timeout    time.Duration
	chromeExec *ChromeExecutor
}

func NewTextCrawler(timeout time.Duration, chromeExecutor *ChromeExecutor) *TextCrawler {
	return &TextCrawler{timeout: timeout, chromeExec: chromeExecutor}
}

func (cl *TextCrawler) FetchOnePage(url string) (string, error) {
	ctx, cancel := cl.chromeExec.NewContext()
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, cl.timeout)
	defer cancel()
	var res string
	var action chromedp.Action = chromedp.Text("html", &res, chromedp.ByQuery)
	if err := crawWithAction(ctx, url, &action); err != nil {
		return "", err
	}
	return res, nil
}

func (cl *TextCrawler) FetchManyPages(urls []string) ([]string, error) {
	ctx, cancel := cl.chromeExec.NewContext()
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Duration(int(cl.timeout)*len(urls)))
	defer cancel()
	res := make([]string, len(urls))

	for idx, url := range urls {
		var action chromedp.Action = chromedp.Text("html", &res[idx], chromedp.ByQuery)
		if err := crawWithAction(ctx, url, &action); err != nil {
			return nil, err
		}
	}

	return res, nil
}
