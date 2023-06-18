package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
)

type Crawler interface {
	FetchOnePage(url string) (string, error)
}

func crawWithAction(ctx context.Context, url string, extractAction *chromedp.Action) error {
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		*extractAction,
	})
	return err
}

type ChromeExecutor struct {
	ExecContext context.Context
	CancelFunc  context.CancelFunc
}

func NewChromeExecutor() *ChromeExecutor {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	execCtx, execCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &ChromeExecutor{
		ExecContext: execCtx,
		CancelFunc:  execCancel,
	}
}

func (c *ChromeExecutor) NewContext() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(c.ExecContext)
}

func (c *ChromeExecutor) ShutDown() {
	c.CancelFunc()
}
