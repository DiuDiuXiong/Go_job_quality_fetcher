package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

type TextCrawler struct {
	ctx     *context.Context
	cancel  *context.CancelFunc
	timeout time.Duration
}

func NewTextCrawler(timeout time.Duration) *TextCrawler {
	ctx, cancel := newContext()
	return &TextCrawler{ctx: ctx, cancel: cancel, timeout: timeout}
}

func (cl *TextCrawler) FetchOnePage(url string) (string, error) {
	ctx, cancel := context.WithTimeout(*cl.ctx, cl.timeout)
	defer cancel()
	var res string
	var action chromedp.Action = chromedp.Text("html", &res, chromedp.ByQuery)
	if err := crawWithAction(ctx, url, &action); err != nil {
		return "", err
	}
	return res, nil
}
