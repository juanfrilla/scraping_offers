package remoterocketship

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sardanioss/httpcloak/client"
)

type RemoteRocketshipScraper struct {
	Session     *client.Client
	ScraperName string
	Logger      *log.Logger
	impersonate string
}

func NewRemoteRocketshipScraper() *RemoteRocketshipScraper {
	impersonate := utils.RandomImpersonation()
	c := client.NewClient(impersonate)
	scraperName := "SimplyHired"
	return &RemoteRocketshipScraper{
		Session:     c,
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
		impersonate: impersonate,
	}
}

func (rrs *RemoteRocketshipScraper) getJobsRequest() (*client.Response, error) {
	ctx := context.Background()
	url := "https://www.remoterocketship.com/?page=1&sort=DateAdded&jobTitle=scraping%2Ccrawling%2Cscraper%2Ccrawler%2Cdata+acquisition&locations=Europe%2CSpain"

	return rrs.Session.Get(ctx, url, nil)

}

func (rrs *RemoteRocketshipScraper) jobInformationRequest(jobURL string) (*client.Response, error) {
	ctx := context.Background()
	return rrs.Session.Get(ctx, jobURL, nil)

}

func (rrs *RemoteRocketshipScraper) parse(responseJSON *RemoteRockJobList) []models.ScrapedJob {
	var records []models.ScrapedJob
	var singleJob RemoteRockJob
	props := responseJSON.Props
	pageProps := props.PageProps
	jobs := pageProps.InitialJobOpenings

	for jobID, job := range jobs {

		rrs.Logger.Printf("[%s] Parsing job posting %d/%d", rrs.ScraperName, jobID+1, len(jobs))

		companyInfo := job.Company
		companyName := strings.ReplaceAll(companyInfo.Name, " ", "-")
		company := utils.NormalizeString(companyName)
		companyImageURL := companyInfo.ProfilePicURL
		slug := job.Slug

		jobURL := fmt.Sprintf("https://www.remoterocketship.com/company/%s/jobs/%s/", companyName, slug)

		if utils.IsForbidden(company, constants.ForbiddenCompanies) {
			continue
		}

		resp, err := rrs.jobInformationRequest(jobURL)
		if err != nil {
			log.Fatal("error fetching job info: %w", err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
		if err != nil {
			log.Fatal("error parsing job HTML: %w", err)
		}

		utils.GetJSONFromHTML(doc, singleJob)

		p := singleJob.Props
		jobInfoProps := p.PageProps

		title := job.RoleTitle
		jobOpening := jobInfoProps.JobOpening
		description := strings.Join([]string{
			jobOpening.RoleDescription,
			jobOpening.RoleRequirements,
			jobOpening.Benefits,
		}, "\n")

		location := job.Location
		datePosted := job.CreatedAt

		keyword := utils.FindKeywordInDescription(description)
		if keyword != "" && !utils.IsForbidden(company, constants.ForbiddenCompanies) && !utils.IsForbidden(title, constants.ForbiddenKeywords) {

			record := models.ScrapedJob{
				Title:           title,
				Company:         company,
				Location:        location,
				URL:             job.URL,
				DatePosted:      utils.FromTimestampToISOFormat(datePosted.UnixMilli()),
				Modality:        "Remote",
				Platform:        rrs.ScraperName,
				KeywordAppeared: keyword,
				LogoURL:         companyImageURL,
			}
			records = append(records, record)
		} else {
			rrs.Logger.Printf("Discarding title %s from company %s", title, company)
			continue
		}
	}

	return records
}

func (rrs *RemoteRocketshipScraper) Scrape() []models.ScrapedJob {
	resp, err := rrs.getJobsRequest()
	if err != nil {
		fmt.Errorf("error getJobsRequest: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Errorf("error parsing HTML: %w", err)
	}
	var rJSON RemoteRockJobList
	if err := utils.GetJSONFromHTML(doc, &rJSON); err != nil {
		fmt.Errorf("error GettingJSONFromHTML: %w", err)
	}

	return rrs.parse(&rJSON)
}

func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
func (rrs *RemoteRocketshipScraper) makeRequest(urlStr string, maxRetries int) (*client.Response, error) {
	retries := 0
	ctx := context.Background()
	for retries < maxRetries {
		rrs.Logger.Printf("Requesting %s with impersonate %s", urlStr, rrs.impersonate)
		resp, err := rrs.Session.Get(ctx, urlStr, nil)
		if err != nil {
			wait := time.Duration(rand.Intn(5)+1) * time.Second
			rrs.Logger.Printf("Request failed (%v), retrying in %s...", err, wait)
			time.Sleep(wait)
			retries++
			rrs.impersonate = utils.RandomImpersonation()
			rrs.Session = client.NewClient(rrs.impersonate)
			continue
		}
		if resp.StatusCode == 429 {
			wait := time.Duration(rand.Intn(5)+1) * time.Second
			rrs.Logger.Printf("Rate limited. Generating new profile, waiting %s seconds...", wait)
			time.Sleep(wait)
			retries++
			rrs.impersonate = utils.RandomImpersonation()
			rrs.Session = client.NewClient(rrs.impersonate)
			continue
		}

		return resp, nil
	}

	rrs.Logger.Printf("Failed to fetch %s after %d retries.", urlStr, maxRetries)
	return nil, fmt.Errorf("max retries exceeded for %s", urlStr)
}
