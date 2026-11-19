package ashby

import "testing"

func TestScrape(t *testing.T) {
	scraper := NewAshbyScraper()

	jobs, err := scraper.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Scraped %d jobs", len(jobs))

	for _, job := range jobs {
		t.Logf("%+v", job)
	}
}
