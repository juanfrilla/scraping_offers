package simplyhired

import (
	"context"
	"fmt"
	"io"
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

type SimplyHiredScraper struct {
	Session     *client.Client
	ScraperName string
	Logger      *log.Logger
	impersonate string
}

func NewSimplyHiredScraper() *SimplyHiredScraper {
	impersonate := utils.RandomImpersonation()
	c := client.NewClient(impersonate)
	scraperName := "SimplyHired"
	return &SimplyHiredScraper{
		Session:     c,
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
		impersonate: impersonate,
	}
}
func (shs *SimplyHiredScraper) Name() string {
	return shs.ScraperName
}

func (shs *SimplyHiredScraper) getJobsRequest(keyword string) (*client.Response, error) {
	ctx := context.Background()
	url := fmt.Sprintf(
		"https://www.simplyhired.es/search?q=%s&l=espa%%C3%%B1a",
		url.QueryEscape(keyword),
	)
	body := strings.NewReader("")

	resp, err := shs.Session.Post(ctx, url, body, nil)
	if err != nil {
		return nil, fmt.Errorf("error in search POST: %w", err)
	}

	return resp, nil
}

func (shs *SimplyHiredScraper) JobInfoRequest(botURL string) (*client.Response, error) {
	fullURL := "https://www.simplyhired.es" + botURL

	ctx := context.Background()

	resp, err := shs.Session.Get(ctx, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error in POST search: %w", err)
	}

	return resp, nil
}

func (shs *SimplyHiredScraper) Parse(jobsData SimplyHiredJobList, seenURLs map[string]bool) []models.ScrapedJob {
	var records []models.ScrapedJob
	props := jobsData.Props
	pageProps := props.PageProps
	jobs := pageProps.Jobs

	for i, j := range jobs {
		shs.Logger.Printf("Processing offer %d/%d", i+1, len(jobs))

		botURL := strings.Split(j.BotURL, "?")[0]
		if seenURLs[botURL] {
			continue
		}
		seenURLs[botURL] = true

		title := j.Title
		company := j.Company
		if utils.IsForbidden(company, constants.ForbiddenCompanies) {
			shs.Logger.Printf("Discarding company %s", company)
			continue
		}
		if utils.IsForbidden(title, constants.ForbiddenKeywords) {
			shs.Logger.Printf("Discarding title %s from company %s", title, company)
			continue
		}

		resp, err := shs.JobInfoRequest(botURL)
		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
		if err != nil {
			log.Fatal("error parsing job HTML: %w", err)
		}

		singleJobPtr, err := utils.GetJSONFromHTML[SimplyHiredJob](doc)
		if err != nil {
			shs.Logger.Printf("Error extracting job JSON: %v", err)
			continue
		}

		singleJob := *singleJobPtr

		jobInfoProps := singleJob.Props
		jobInfoPageProps := jobInfoProps.PageProps
		description := jobInfoPageProps.JobDescriptionHTML
		if description == "" {
			description = "N/A"
		}

		keyword := utils.FindKeywordInDescription(description)

		if keyword == "" {
			shs.Logger.Printf("Discarding offer %s from company %s", title, company)
			continue
		}

		record := models.ScrapedJob{
			Title:           title,
			Company:         company,
			Location:        j.Location,
			URL:             "https://www.simplyhired.es" + botURL,
			DatePosted:      utils.FromTimestampToTime(j.DateOnIndeed),
			Modality:        utils.DetermineModality(title, description),
			Platform:        shs.ScraperName,
			KeywordAppeared: keyword,
			LogoURL:         jobInfoPageProps.EmployerSquareLogoURL,
			ExternalID:      utils.Slugify(company + ":" + j.JobKey),
		}

		records = append(records, record)
	}

	return records
}

func (shs *SimplyHiredScraper) Scrape() ([]models.ScrapedJob, error) {
	var allJobs []models.ScrapedJob

	seenURLs := make(map[string]bool)

	for _, kw := range constants.SearchKeywords {
		shs.Logger.Printf("Scraping word: %s", kw)

		resp, err := shs.getJobsRequest(kw)
		if err != nil {
			return nil, fmt.Errorf("error getJobsRequest para %s: %w", kw, err)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("error parsing HTML para %s: %w", kw, err)
		}

		rJSON, err := utils.GetJSONFromHTML[SimplyHiredJobList](doc)
		if err != nil {
			return nil, fmt.Errorf("error GetJSONFromHTML para %s: %w", kw, err)
		}

		if len(rJSON.Props.PageProps.Jobs) > 0 {
			jobs := shs.Parse(*rJSON, seenURLs)
			allJobs = append(allJobs, jobs...)

			shs.Logger.Printf("Found %d jobs for %s", len(jobs), kw)
		}
	}

	return allJobs, nil
}
