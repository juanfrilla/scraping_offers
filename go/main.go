package main

import (
	"encoding/json"
	"fmt"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/db"
	"scraping_offers/go/models"
	"scraping_offers/go/scrapers/remoterocketship"
	"scraping_offers/go/scrapers/simplyhired"
	"scraping_offers/go/utils"

	"sync"

	"github.com/joho/godotenv"
)

type Scraper interface {
	Name() string
	Scrape() ([]models.ScrapedJob, error)
}

func main() {
	if err := godotenv.Load(constants.EnvFile); err != nil {
		fmt.Println("No .env file found")
	}

	//TODO estas carpetas ponerlas en constants
	database, err := db.InitDB(constants.DataDB)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}
	defer database.Close()

	os.MkdirAll("data", 0755)
	os.MkdirAll("debug", 0755)
	scrapers := []Scraper{
		//		linkedin.NewLinkedinScraper(),
		//		infojobs.NewInfojobsScraper(),
		//		ashby.NewAshbyScraper(),
		remoterocketship.NewRemoteRocketshipScraper(),
		// remoteok.NewRemoteOKScraper(),
		// empleate.NewEmpleateScraper(),
	}

	if utils.IsLocal() {
		scrapers = append(scrapers, simplyhired.NewSimplyHiredScraper())
	}

	resultsChan := make(chan []models.ScrapedJob, len(scrapers))
	errorsChan := make(chan error, len(scrapers))
	var wg sync.WaitGroup
	for _, scraper := range scrapers {
		wg.Add(1)
		go func(s Scraper) {
			defer wg.Done()
			fmt.Println("Starting Scraping", s.Name())

			data, err := s.Scrape()

			if err != nil {
				fmt.Println("❌ Scraping failed", s.Name(), err)
				errorsChan <- fmt.Errorf("%s: %w", s.Name(), err)
				return
			}

			fmt.Println("Finished Scraping", s.Name())

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

	fmt.Println("Working directory:")
	wd, _ := os.Getwd()
	fmt.Println(wd)

	fmt.Println("DataFile:", constants.DataFile)

	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Println("Error creating data directory:", err)
		return
	}

	file, err := os.Create(constants.DataFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	if allJobs == nil {
		allJobs = []models.ScrapedJob{}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allJobs); err != nil {
		fmt.Println("Error writing JSON:", err)
		os.Exit(1)
	}

	err = utils.UpdateLastScraped()
	if err != nil {
		fmt.Println("Error updating last scraped:", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Exported %d jobs to %s\n", len(allJobs), constants.DataFile)
}
