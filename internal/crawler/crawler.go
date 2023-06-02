package crawler

type Crawler interface {
	FetchOnePage(url string) (string, error)
}
