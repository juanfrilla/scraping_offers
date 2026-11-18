package ashby

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/sardanioss/httpcloak/client"
)

type AshbyScraper struct {
	Session     *client.Client
	ScraperName string
	Logger      *log.Logger
	impersonate string
}

func NewAshbyScraper() *AshbyScraper {
	scraperName := "Ashby"
	impersonate := utils.RandomImpersonation()
	c := client.NewClient(impersonate)
	return &AshbyScraper{
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
		Session:     c,
		impersonate: impersonate,
	}
}

func (ls *AshbyScraper) Name() string {
	return ls.ScraperName
}
func (s *AshbyScraper) githubAshbyRequest() (*client.Response, error) {
	ctx := context.Background()
	url := "https://github.com/george-kozlitin/ats-jobs-directory/blob/main/directory-ashby.json"
	return s.Session.Get(ctx, url, nil)

}

func (s *AshbyScraper) retrieveCompanyJobs(slug string) (*AshbyJobsList, error) {
	endpoint := fmt.Sprintf(
		"https://api.ashbyhq.com/posting-api/job-board/%s",
		slug,
	)

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("includeCompensation", "true")
	parsedURL.RawQuery = query.Encode()

	headers := map[string][]string{
		"accept": {
			"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		},
		"accept-language": {"es-ES,es;q=0.9,en;q=0.8"},
		"cache-control":   {"no-cache"},
		"pragma":          {"no-cache"},
	}

	resp, err := s.Session.Get(
		context.Background(),
		parsedURL.String(),
		headers,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := resp.Text()
		return nil, fmt.Errorf(
			"HTTP %d: %s",
			resp.StatusCode,
			body,
		)
	}

	var data AshbyJobsList

	if err := resp.JSON(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
func (s *AshbyScraper) retrieveAshbyData() (*ghRawdata, error) {
	resp, err := s.githubAshbyRequest()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))

	embeddedData, err := utils.GetJSONFromHTML[ghJobsPreData](doc)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	// Buscar:
	// <script data-target="react-app.embeddedData">

	payload := embeddedData.Payload

	route := payload.CodeViewBlobLayoutRouteStyledBlob //.(map[string]interface{})

	rawLines := route.RawLines

	rawJSON := rawLines[0]

	var data ghRawdata

	if err := json.Unmarshal([]byte(rawJSON), &data); err != nil {
		return nil, err
	}

	if err != nil {
		s.Logger.Printf("Error parsing JSON after fix: %v", err)
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	return &data, nil
}

func (s *AshbyScraper) Scrape() ([]models.ScrapedJob, error) {
	ashbyData, err := s.retrieveAshbyData()
	if err != nil {
		return nil, err
	}

	companies := ashbyData.Companies

	collected := make([]models.ScrapedJob, 0)

	for index, company := range companies {
		s.Logger.Printf(
			"Processing %d/%d %s\n",
			index+1,
			len(companies),
			company.Company,
		)
		slug := company.Slug

		if slug == "" {
			continue
		}

		companyData, err := s.retrieveCompanyJobs(company.Slug)
		if err != nil {
			s.Logger.Printf(
				"  ! Error en %s: %v\n",
				company.Slug,
				err,
			)
			continue
		}

		for index, job := range companyData.Jobs {
			s.Logger.Printf(
				"Processing %s - %d/%d\n",
				company.Company,
				index+1,
				len(companyData.Jobs),
			)
			var companyStr string
			if company.Company == "" {
				companyStr = slug
			} else {
				companyStr = company.Company

			}

			location := job.Address.PostalAddress.AddressCountry

			description := job.DescriptionPlain

			isRemote := job.IsRemote == nil || *job.IsRemote
			keyword := utils.FindKeywordInDescription(description)

			if job.JobURL != "" &&
				keyword != "" &&
				isRemote &&
				!utils.IsForbidden(job.Title, constants.ForbiddenKeywords) && !utils.IsForbidden(location, constants.ForbiddenCountries) {

				var modality string

				if job.IsRemote != nil && *job.IsRemote {
					modality = "Remote"
				} else {
					modality = "Unknown"
				}

				s.Logger.Printf(
					"Appending %s\n",
					job.Title,
				)

				collected = append(
					collected,
					models.ScrapedJob{
						Title:           job.Title,
						Company:         companyStr,
						Location:        location,
						URL:             job.JobURL,
						DatePosted:      utils.FromTimestampToISOFormat(job.PublishedAt.UnixMilli()),
						Modality:        modality,
						Platform:        "Ashby",
						KeywordAppeared: keyword,
						LogoURL:         "",
					},
				)
			}
		}
	}
	return collected, nil
}
