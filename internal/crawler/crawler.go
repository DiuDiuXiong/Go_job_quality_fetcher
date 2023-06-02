package crawler

type Crawler interface {
	FetchOnePage(url string) (string, error)
}

type BasicCrawler struct {
}

func (bc *BasicCrawler) FetchOnePage(url string) (string, error) {
	return "", nil
}
