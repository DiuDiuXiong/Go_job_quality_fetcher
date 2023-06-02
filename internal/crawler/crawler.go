package crawler

type crawler interface {
	FetchOnePage(url string) (string, error)
}

type BasicCrawler struct {
}
