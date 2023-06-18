package fetcher

type FetchCriteria struct {
	JobTitle    string
	Duration    int // number of days to look back
	Description string
	Country     string
}

type Fetcher interface {
	FetchContents(criteria *FetchCriteria) ([]string, error)              // get all contents based on search criteria
	FetchCriteriaToUrl(fc *FetchCriteria, pageNumber int) (string, error) // transfer to query urls (first three page)
}
