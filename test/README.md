# Test Structure

## `/internal/fetcher/`
Contains tests required for fetcher. This will also generate files under tmp for online CI check.
- `/tmp`: storage for fetched html files

## `/internal/crawler`
- [`html_test.go`](internal/crawler/html_test.go)
This one ensures crawler is not blocked by target websites, also ensures the file are generated for online CI check.
- [`text_test.go`](internal/crawler/text_test.go)
This one mainly ensures that the `ExtractJobDescriptionUrls` are returning a right format