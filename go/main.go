package main

import (
	"encoding/json"
	"fmt"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/scrapers/infojobs"
	"scraping_offers/go/scrapers/linkedin"
	"scraping_offers/go/scrapers/remoteok"
	"scraping_offers/go/scrapers/remoterocketship"
	"scraping_offers/go/scrapers/simplyhired"

	"sync"
)

type Scraper interface {
	Name() string
	Scrape() ([]models.ScrapedJob, error)
}

func main() {
	os.MkdirAll("data", 0755)
	os.MkdirAll("debug", 0755)
	scrapers := []Scraper{
		linkedin.NewLinkedinScraper(),
		infojobs.NewInfojobsScraper(),
		simplyhired.NewSimplyHiredScraper(),
		remoteok.NewRemoteOKScraper(),
		remoterocketship.NewRemoteRocketshipScraper(),
	}

	resultsChan := make(chan []models.ScrapedJob)
	errorsChan := make(chan error)
	var wg sync.WaitGroup
	for _, scraper := range scrapers {
		wg.Add(1)
		go func(s Scraper) {
			defer wg.Done()

			data, err := s.Scrape()
			if err != nil {
				errorsChan <- fmt.Errorf("%s: %w", s.Name(), err)
				return
			}

			resultsChan <- data
		}(scraper)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()
	var allJobs []models.ScrapedJob

	for res := range resultsChan {
		allJobs = append(allJobs, res...)
	}

	for err := range errorsChan {
		fmt.Println("❌ Error in scraper:", err)
	}

	file, err := os.Create(constants.DataFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allJobs); err != nil {
		fmt.Println("Error writing JSON:", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Exported %d jobs to %s\n", len(allJobs), constants.DataFile)
}
