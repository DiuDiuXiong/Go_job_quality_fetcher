package crawler

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"time"
)

func generateCrawlerFunction(actionGenerator func(*string) chromedp.Action, guiEnable bool) func(string) (string, error) {
	return func(url string) (string, error) {

		headers := map[string]interface{}{
			"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
			"sec-ch-ua":          `"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"macOS"`,
		}

		opts := chromedp.DefaultExecAllocatorOptions[:]
		if guiEnable {
			opts = append(opts, chromedp.Flag("headless", false))
		}
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		ctx, cancel = chromedp.NewContext(ctx)
		defer cancel()

		var res string
		action := actionGenerator(&res)

		err := chromedp.Run(ctx, chromedp.Tasks{
			network.Enable(), // this allows further network event (listening/ set header ...)
			network.SetExtraHTTPHeaders(headers),
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

/*
	To bypass cloudflare, this will also work:
		//opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//	chromedp.Flag("headless", false),
		//)
	However, the issue is aws lambda/ec2... don't have a gui. So this will still cause a headless header in http request, which will cause block

	The header is copy from header with non-headless (hence no block). By code below:

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Enable the network
	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		log.Fatalf("Could not enable network: %v", err)
	}

	// Set up a listener for the requestWillBeSent event
	requests := make(map[network.RequestID]*network.Request)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			requests[ev.RequestID] = ev.Request
		}
	})

	// Visit the URL
	if err := chromedp.Run(ctx, chromedp.Navigate(`https://au.indeed.com/`)); err != nil {
		log.Fatalf("Could not navigate to url: %v", err)
	}

	// Print all requests
	for _, req := range requests {
		log.Printf("Request to %s with headers %v", req.URL, req.Headers)
	}

	// Get window size and user agent
	var size windowSize
	var userAgent string
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`({width: window.innerWidth, height: window.innerHeight})`, &size),
		chromedp.Evaluate(`navigator.userAgent`, &userAgent),
	)
	if err != nil {
		log.Fatalf("Could not get window size or user agent: %v", err)
	}

	log.Printf("Window size: %+v", size)
	log.Printf("User agent: %s", userAgent)
*/
