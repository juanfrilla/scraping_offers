package infojobs

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	targetURL := "https://www.infojobs.net/jobsearch/search-results/list.xhtml"

	formData := url.Values{}
	formData.Set("palabra", keyword)
	formData.Set("normalizedJobTitleId", "")
	formData.Set("of_provincia", "0")
	formData.Set("canal", "0")
	formData.Set("origen_busqueda", "0")
	formData.Set("origen_accion", "0")
	formData.Set("vieneUrlExecutive", "false")

	body := strings.NewReader(formData.Encode())

	resp, err := c.Post(ctx, targetURL, body, nil)
	if err != nil {
		return nil, fmt.Errorf("error in POST search: %w", err)
	}

	return resp, nil
}

func (ijs *InfojobsScraper) Scrape() ([]models.ScrapedJob, error) {
	var records []models.ScrapedJob
	var rawData Root

	for _, kw := range constants.SearchKeywords {
		ijs.Logger.Printf("Searching for: %s", kw)
		resp, err := ijs.InfojobsJobSearchRequest(kw)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			ijs.Logger.Printf("Bad status for %s: %d", kw, resp.StatusCode)
			continue
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)

		// html, err := doc.Html()

		// htmlFilename := fmt.Sprintf("debug_%s.html", kw)
		// err = os.WriteFile(htmlFilename, []byte(html), 0644)

		utils.GetJSONFromHTML(doc, &rawData)

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
					DatePosted:      offer.PublishedAt,
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
