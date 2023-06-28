package main

import (
	"fmt"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/crawler"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/duplicateremover"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/extractor"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/fetcher"
	"github.com/DiuDiuXiong/Go_job_quality_fetcher/internal/summarisation"
	"io/ioutil"
	"path"
	"strconv"
	"sync"
	"time"
)

func main() {

	// 01 Get Job Title
	fmt.Println("Insert job title:")
	var jobTitle string
	fmt.Scan(&jobTitle)
	if len(jobTitle) == 0 {
		fmt.Println("No Job Title Provided")
		return
	}

	// 02 Start Crawler
	duration := 30 * time.Second
	indeed := fetcher.NewIndeedFetcher(duration, crawler.NewChromeExecutor())
	seek := fetcher.NewSeekFetcher(duration, crawler.NewChromeExecutor())
	linkedin := fetcher.NewLinkedinFetcher(duration, crawler.NewChromeExecutor())

	// 03 Start Crawling Jobs
	query := &fetcher.FetchCriteria{
		JobTitle:    jobTitle,
		Duration:    30,
		Description: "-",
		Country:     "Australia",
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// 03-2 Show Progress Bar
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

	// 04 remove duplicates
	fmt.Println("Crawling finish. Start remove duplicates...")
	err := duplicateremover.RemoveDuplicates(
		[]string{path.Join("store", "linkedin"), path.Join("store", "seek"), path.Join("store", "indeed")},
		path.Join("store", "total"), 0.55)

	if err != nil {
		fmt.Println(fmt.Errorf("failed to remove duplicates"))
		fmt.Println(fmt.Errorf(err.Error()))
		return
	}

	// 05 Query Open AI for Summary
	openai := summarisation.NewOpenAIClient()

	// extract key job qualities
	err = openai.GetTechnicalQualitiesForJobs(path.Join("store", "total"), path.Join("store", "result.txt"))
	if err != nil {
		fmt.Println(fmt.Errorf("failed to remove extract job qualities"))
		fmt.Println(fmt.Errorf(err.Error()))
		return
	}

	// get suggestion

	fmt.Println("Successfully extract job qualities, getting summary from those qualities")
	dt, _ := ioutil.ReadFile(path.Join("store", "result.txt"))
	keyJobQualities, err := openai.GetTechnicalQualitiesCountSummary(dt)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get most occured qualities"))
		fmt.Println(fmt.Errorf(err.Error()))
		return
	}

	fmt.Printf("\n\n\n\n\n\n\n\nThe Key Job Qualities for %s are", jobTitle)
	fmt.Println(keyJobQualities)

}
