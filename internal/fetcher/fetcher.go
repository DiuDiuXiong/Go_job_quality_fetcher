package fetcher

type fether interface {
	FetchContents(jobTitle string, criteria map[string]interface{}) []string
}

type BasicFetcher struct {
}
