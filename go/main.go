package main

import (
	"encoding/json"
	"fmt"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/db"
	"scraping_offers/go/models"
	"scraping_offers/go/scrapers/ashby"
	"scraping_offers/go/scrapers/empleate"
	"scraping_offers/go/scrapers/infojobs"
	"scraping_offers/go/scrapers/linkedin"
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
	if err := constants.InitPaths(); err != nil {
		fmt.Println("Error initializing paths:", err)
		os.Exit(1)
	}
	if err := godotenv.Load(constants.EnvFile); err != nil {
		fmt.Println("No .env file found")
	}

	database, err := db.InitDB(constants.DataDB)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}
	defer database.Close()

	os.MkdirAll("data", 0755)
	os.MkdirAll("debug", 0755)
	scrapers := []Scraper{
		linkedin.NewLinkedinScraper(),
		infojobs.NewInfojobsScraper(),
		ashby.NewAshbyScraper(),
		remoterocketship.NewRemoteRocketshipScraper(),
		// remoteok.NewRemoteOKScraper(),
		empleate.NewEmpleateScraper(),
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

	if err := db.SaveJobs(database, allJobs); err != nil {
		fmt.Println("Error saving jobs to database:", err)
		os.Exit(1)
	}

	var count int
	var deleted int64
	count, err = db.GetJobsCount(database)
	if err != nil {
		fmt.Println("Error checking jobs count:", err)
		os.Exit(1)
	}

	fmt.Println("Jobs count:", count)

	file, err = os.Create(constants.DataFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(allJobs); err != nil {
		file.Close()
		fmt.Println("Error writing JSON:", err)
		os.Exit(1)
	}

	if err := file.Close(); err != nil {
		fmt.Println("Error closing JSON file:", err)
		os.Exit(1)
	}

	if err := utils.UpdateLastScraped(); err != nil {
		fmt.Println("Error updating last scraped:", err)
		os.Exit(1)
	}
	if deleted, err = db.DeleteJobsNotUpdatedToday(database); err != nil {
		fmt.Println("Error deleting jobs:", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Saved %d jobs to SQLite\n", len(allJobs))
	fmt.Printf("X Jobs deleted: %d\n", deleted)
	fmt.Printf("📊 Jobs actually in database: %d\n", count)
}
