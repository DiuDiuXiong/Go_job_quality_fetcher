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

func newContext() (*context.Context, *context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(ctx)
	return &ctx, &cancel
}
