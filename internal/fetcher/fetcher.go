package fetcher

type fether interface {
	FetchContents(jobTitle string, criteria map[string]interface{}) []string
}

type crawler interface {
	FetchOnePage(url string) (string, error)
}

type BasicCrawler struct {
}

type BasicFetcher struct {
}
