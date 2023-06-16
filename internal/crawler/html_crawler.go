package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

type HTMLCrawler struct {
	ctx     *context.Context
	cancel  *context.CancelFunc
	timeout time.Duration
}

func NewHTMLCrawler(timeout time.Duration) *HTMLCrawler {
	ctx, cancel := newContext()
	return &HTMLCrawler{ctx: ctx, cancel: cancel, timeout: timeout}
}

func (cl *HTMLCrawler) FetchOnePage(url string) (string, error) {
	ctx, cancel := context.WithTimeout(*cl.ctx, cl.timeout)
	defer cancel()
	var res string
	var action chromedp.Action = chromedp.OuterHTML("html", &res, chromedp.ByQuery)
	if err := crawWithAction(ctx, url, &action); err != nil {
		return "", err
	}
	return res, nil

}
