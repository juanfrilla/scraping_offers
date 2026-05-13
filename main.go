package main

import (
	"encoding/json"
	"fmt"
	"os"
	"scraping_offers/constants"
	"scraping_offers/models"

	//"scraping_offers/scrapers/infojobs"

	"scraping_offers/scrapers/infojobs"
	"scraping_offers/scrapers/linkedin"
	"scraping_offers/scrapers/remoteok"
	"scraping_offers/scrapers/remoterocketship"
	"scraping_offers/scrapers/simplyhired"
	// "scraping_offers/scrapers/remoteok"
	// "scraping_offers/scrapers/remoterocketship"
	// "scraping_offers/scrapers/simplyhired"
	"sync"
)

// type Scraper interface {
// 	Scrape() ([]models.ScrapedJob, error)
// } // arreglarla mañana

type Scraper interface {
	Scrape() []models.ScrapedJob
}

func main() {
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

	// Lanzar scrapers en paralelo
	for _, scraper := range scrapers {
		wg.Add(1)
		go func(s Scraper) {
			defer wg.Done()
			data := s.Scrape()
			resultsChan <- data
		}(scraper)
	}

	// Cerrar canales cuando terminen todos
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	// Recoger resultados
	var allJobs []models.ScrapedJob

	for res := range resultsChan {
		allJobs = append(allJobs, res...)
	}

	// Manejar errores (si quieres)
	for err := range errorsChan {
		fmt.Println("❌ Error en scraper:", err)
	}

	// Guardar JSON
	os.MkdirAll("data", 0755)

	file, err := os.Create(constants.DataFile)
	if err != nil {
		fmt.Println("Error creando archivo:", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allJobs); err != nil {
		fmt.Println("Error escribiendo JSON:", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Exportados %d jobs a %s\n", len(allJobs), constants.DataFile)
}
