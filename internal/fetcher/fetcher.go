package fetcher

type Fether interface {
	FetchContents(jobTitle string, criteria map[string]interface{}) []string
}

type BasicFetcher struct {
}
