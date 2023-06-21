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
	ExtractJobDescriptionUrls(html *string) []string                      // extract link to actual jobs based on a search page results
	ExtractTotalJobCounts(html *string) int                               // extract number job gets based on a sort
	FetchTargetJobUrls(fc *FetchCriteria) ([]string, error)               // wrapper for extractJobDescriptionUrls to define stop condition (we might interested in jobs more than one page)
}
