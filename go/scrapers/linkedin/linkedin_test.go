package linkedin

import (
	"fmt"
	"log"
	"testing"
)

func TestScrape(t *testing.T) {
	scraper := NewLinkedinScraper()

	if scraper == nil {
		log.Fatal("Error creating scraper")
	}

	jobs, err := scraper.Scrape()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Scraped %d jobs\n", len(jobs))

	for _, job := range jobs {
		fmt.Printf("%+v\n", job)
	}
}
