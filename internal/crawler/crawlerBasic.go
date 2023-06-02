package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

func generateCrawlerFunction(actionGenerator func(string) chromedp.Action) func(string) (string, error) {
	return func(url string) (string, error) {
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[:]...)
		defer cancel()
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		ctx, cancel = chromedp.NewContext(ctx)
		defer cancel()

		var res string
		action := actionGenerator(res)

		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(url),
			action,
		})

		if err != nil {
			// log error to some logging service
			return "", err
		}

		return res, nil
	}
}

/*
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	:make the window UI present when chromedp.Run

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[:]...)
	defer cancel()

	:allocate space for chrome instance/ set up environment variable for chrome instance running

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	This can be called multiple time without repeated NewExecAllocator, it sets up an isolated context

	Run runs a series of actions
*/
