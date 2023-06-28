package main

import (
	"fmt"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/summarisation"
	"io/ioutil"
	"path"
)

func main() {
	/*
		duration := 30 * time.Second
		chromeExecutor_seek := crawler.NewChromeExecutor()
		chromeExecutor_indeed := crawler.NewChromeExecutor()
		chromeExecutor_linkedin := crawler.NewChromeExecutor()
		//htmlFetcher := crawler.NewHTMLCrawler(duration, chromeExecutor)
		//textFetcher := crawler.NewTextCrawler(duration, chromeExecutor)
		//html, _ := htmlFetcher.FetchOnePage("https://au.linkedin.com/jobs/api/seeMoreJobPostings/search?keywords=SRE&location=Australia&start=0")
		//p, _ := textFetcher.FetchOnePage("https://au.indeed.com/viewjob?jk=1dc3f68b05fb0f7e&tk=1h3mhffdkk4kj800&from=hp&advn=846196049671291&adid=403884505&ad=-6NYlbfkN0CuwhhmBKfyXQphaAl4pby4iVb9T3cBPzknJFM75qC5D52jZLke-URFsTRqZpt7DSv_Ln5SHKWQggzn_stZ6wa5mxMCXpopm_cPDoUV1LPThZOedmlDrZ2b566nN-44C2bolEo9kB1puNbcFj9JWRU2ie8J2eeAS28Q80Hrqyt6NdFvIrTuEJSty5qclVDXI0dSlUWADCfZ3bJ0Ot290Il6N8sI05fwMxaxvaQswRiYRITVdVpx1y7UydNsxZo9snR0wms0bHg0w3CfTZjO9WCHTXdZOM_E9FYIcfhRaQMmW1coaovf5MoF19mFBc_NXTqT9J9-yjbZ2xv7QNGvxhKJcj1euMphT6dNbvfoo-OS0kTuZXMvqK1FUwgkweETiyDj1jX_vKzkG2nO-n8D3pxZyi_4oQHtNrnMt0Br_69v4Yt7amxo6pnAjBabpFKUFMi91sG9tzwfoZj8fIOrdDKU76Jy7lEIOTEhTWQld16W5Q%3D%3D&pub=4a1b367933fd867b19b072952f68dceb&xkcb=SoBs-_M3O2qyKcwt650JbzkdCdPP&vjs=3")
		//p, _ := ioutil.ReadFile(path.Join("Linkedin_0.txt"))
		//pp := string(p)
		//fmt.Println(extractor.ExtractLinkedinJobDescription(&pp))

		query := &fetcher.FetchCriteria{
			JobTitle:    "SRE",
			Duration:    50,
			Description: "-",
			Country:     "Australia",
		}
		indeed := fetcher.NewIndeedFetcher(duration, chromeExecutor_indeed)
		seek := fetcher.NewSeekFetcher(duration, chromeExecutor_seek)
		linkedin := fetcher.NewLinkedinFetcher(duration, chromeExecutor_linkedin)

		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			resIndeed, _ := indeed.FetchContents(query)
			for idx, txt := range resIndeed {
				txtShort := extractor.ExtractIndeedJobDescription(&txt)
				ioutil.WriteFile(path.Join("store", "indeed", strconv.Itoa(idx)+".txt"), []byte(txtShort), 0644)
			}
			wg.Done()
		}()
		go func() {
			resSeek, _ := seek.FetchContents(query)
			for idx, txt := range resSeek {
				txtShort := extractor.ExtractSeekJobDescription(&txt)
				ioutil.WriteFile(path.Join("store", "seek", strconv.Itoa(idx)+".txt"), []byte(txtShort), 0644)
			}
			wg.Done()
		}()
		go func() {
			resLinkedin, _ := linkedin.FetchContents(query)
			for idx, txt := range resLinkedin {
				txtShort := extractor.ExtractLinkedinJobDescription(&txt)
				ioutil.WriteFile(path.Join("store", "linkedin", strconv.Itoa(idx)+".txt"), []byte(txtShort), 0644)
			}
			wg.Done()
		}()

		wg.Wait()
	*/
	/*
		duplicateremover.RemoveDuplicates(
			[]string{path.Join("store", "linkedin"), path.Join("store", "seek"), path.Join("store", "indeed")},
			path.Join("store", "total"),
			0.55)*/

	openai := summarisation.NewOpenAIClient()
	dt, _ := ioutil.ReadFile(path.Join("store", "result.txt"))
	fmt.Println(openai.GetTechnicalQualitiesCountSummary(dt))
}

func ContainsDuplicates(nums []string) bool {
	encountered := make(map[string]bool)

	for _, num := range nums {
		if encountered[num] {
			return true
		}
		encountered[num] = true
	}

	return false
}
