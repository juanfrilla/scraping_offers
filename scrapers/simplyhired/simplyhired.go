package simplyhired

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"scraping_offers/constants"
	"scraping_offers/models"
	"scraping_offers/utils"
	"strings"
	"time"

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

func (shs *SimplyHiredScraper) getJobsRequest(keyword string) (*client.Response, error) {
	ctx := context.Background()
	url := fmt.Sprintf("https://www.simplyhired.es/search?q=%s&l=espa%C3%B1a", keyword)
	body := strings.NewReader("")

	resp, err := shs.Session.Post(ctx, url, body, nil)
	if err != nil {
		return nil, fmt.Errorf("error en POST búsqueda: %w", err)
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
		var singleJob SimplyHiredJob
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

		utils.GetJSONFromHTML(doc, &singleJob)

		jobInfoProps := singleJob.Props
		jobInfoPageProps := jobInfoProps.PageProps
		description := jobInfoPageProps.JobDescriptionHTML
		if description == "" {
			description = "N/A"
		}

		keyword := utils.FindKeywordInDescription(description)
		if keyword != "" {
			record := models.ScrapedJob{
				Title:           title,
				Company:         company,
				Location:        j.Location,
				URL:             "https://www.simplyhired.es" + botURL,
				DatePosted:      utils.FromTimestampToISOFormat(j.DateOnIndeed),
				Modality:        utils.DetermineModality(title, description),
				Platform:        shs.ScraperName,
				KeywordAppeared: keyword,
				LogoURL:         jobInfoPageProps.EmployerSquareLogoURL,
			}
			records = append(records, record)
		} else {
			shs.Logger.Printf("Discarding title %s from company %s", title, company)
			continue
		}
	}

	return records
}

func (shs *SimplyHiredScraper) Scrape() []models.ScrapedJob {
	var allJobs []models.ScrapedJob
	seenURLs := make(map[string]bool)
	for _, kw := range constants.SearchKeywords {
		shs.Logger.Printf("Scraping word: %s", kw)
		resp, err := shs.getJobsRequest(kw)
		if err != nil {
			shs.Logger.Printf("Error getJobsRequest para %s: %v", kw, err)
			continue
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			shs.Logger.Printf("Error parsing HTML para %s: %v", kw, err)
			continue
		}
		var rJSON SimplyHiredJobList
		if err := utils.GetJSONFromHTML(doc, &rJSON); err != nil {
			shs.Logger.Printf("Error GettingJSONFromHTML para %s: %v", kw, err)
			continue
		}
		if len(rJSON.Props.PageProps.Jobs) > 0 {
			jobs := shs.Parse(rJSON, seenURLs)
			allJobs = append(allJobs, jobs...)
			shs.Logger.Printf("Found %d jobs for %s", len(jobs), kw)
		}
	}

	return allJobs
}
func (shs *SimplyHiredScraper) makeRequest(urlStr string, maxRetries int) (*client.Response, error) {
	retries := 0
	ctx := context.Background()
	for retries < maxRetries {
		shs.Logger.Printf("Requesting %s with impersonate %s", urlStr, shs.impersonate)
		resp, err := shs.Session.Get(ctx, urlStr, nil)
		if err != nil {
			wait := time.Duration(rand.Intn(5)+1) * time.Second
			shs.Logger.Printf("Request failed (%v), retrying in %s...", err, wait)
			time.Sleep(wait)
			retries++
			shs.impersonate = utils.RandomImpersonation()
			shs.Session = client.NewClient(shs.impersonate)
			continue
		}
		if resp.StatusCode == 429 {
			wait := time.Duration(rand.Intn(5)+1) * time.Second
			shs.Logger.Printf("Rate limited. Generating new profile, waiting %s seconds...", wait)
			time.Sleep(wait)
			retries++
			shs.impersonate = utils.RandomImpersonation()
			shs.Session = client.NewClient(shs.impersonate)
			continue
		}

		return resp, nil
	}

	shs.Logger.Printf("Failed to fetch %s after %d retries.", urlStr, maxRetries)
	return nil, fmt.Errorf("max retries exceeded for %s", urlStr)
}
