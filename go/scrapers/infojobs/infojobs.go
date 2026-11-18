package infojobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"

	"github.com/sardanioss/httpcloak/client"
)

type InfojobsScraper struct {
	ScraperName string
	Logger      *log.Logger
	impersonate string
}

func NewInfojobsScraper() *InfojobsScraper {
	impersonate := utils.RandomImpersonation()
	scraperName := "Infojobs"
	return &InfojobsScraper{
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
		impersonate: impersonate,
	}
}

func (ijs *InfojobsScraper) Name() string {
	return ijs.ScraperName
}

func (ijs *InfojobsScraper) InfojobsJobSearchRequest(keyword string) (*client.Response, error) {
	ctx := context.Background()
	ijs.Logger.Printf("Using %s", ijs.impersonate)
	c := client.NewClient(ijs.impersonate)
	targetURL := "https://www.infojobs.net/webapp/offers/search?keyword=%22scraping%22&searchByType=empty&referer=search-new"

	resp, err := c.Get(ctx, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error in POST search: %w", err)
	}

	return resp, nil
}

func (ijs *InfojobsScraper) Scrape() ([]models.ScrapedJob, error) {
	var records []models.ScrapedJob
	var rawData IjData
	for _, kw := range constants.SearchKeywords {
		ijs.Logger.Printf("Searching for: %s", kw)
		resp, err := ijs.InfojobsJobSearchRequest(kw)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			ijs.Logger.Printf("Bad status for %s: %d", kw, resp.StatusCode)
			resp.Body.Close() // Importante cerrar el body
			continue
		}

		resp.Body.Close() // Cerramos el body original ya que tenemos los datos en memoria
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
			panic(fmt.Errorf("error decodificando JSON: %w", err))
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing JSON from HTML: %w", err)
		}

		if err != nil {
			return nil, fmt.Errorf("error parsing HTML: %w", err)
		}

		ijs.Logger.Printf("Successfully scraped %s", kw)
		for _, offer := range rawData.Offers {
			title := offer.Title
			description := offer.Description
			company := utils.NormalizeString(offer.CompanyName)
			keyword := utils.FindKeywordInDescription(description)
			if keyword != "" && !utils.IsForbidden(company, constants.ForbiddenCompanies) && !utils.IsForbidden(title, constants.ForbiddenKeywords) {
				job := models.ScrapedJob{
					Title:           title,
					Company:         company,
					Location:        offer.City,
					URL:             offer.Link,
					DatePosted:      utils.FromTimestampToISOFormat(offer.PublishedAt.UnixMilli()),
					Modality:        utils.DetermineModality(offer.Title, offer.Description),
					Platform:        ijs.ScraperName,
					KeywordAppeared: kw,
					LogoURL:         offer.CompanyLogo,
				}

				records = append(records, job)
			}
		}
	}

	return records, nil
}
