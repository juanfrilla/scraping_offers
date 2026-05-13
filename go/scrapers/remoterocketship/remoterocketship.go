package remoterocketship

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"
	"strings"

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
	scraperName := "RemoteRocketship"
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
			rrs.Logger.Printf("Discarding company %s", company)
			continue
		}

		var description string

		if utils.IsLocal() {
			resp, err := rrs.jobInformationRequest(jobURL)
			if err != nil {
				log.Printf("error fetching job info: %v", err)
				continue
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("error reading body: %v", err)
				continue
			}

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
			if err != nil {
				log.Printf("error parsing job HTML: %v", err)
				continue
			}

			utils.GetJSONFromHTML(doc, &singleJob)

			p := singleJob.Props
			jobInfoProps := p.PageProps
			jobOpening := jobInfoProps.JobOpening

			description = strings.Join([]string{
				jobOpening.RoleDescription,
				jobOpening.RoleRequirements,
				jobOpening.Benefits,
			}, "\n")
		} else {
			description = ""
		}
		title := job.RoleTitle
		location := job.Location
		datePosted := job.CreatedAt
		var keyword string
		if utils.IsLocal() {

			keyword = utils.FindKeywordInDescription(description)
		} else {
			keyword = ""
		}

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
	}

	return records
}

// func (rrs *RemoteRocketshipScraper) parse(responseJSON *RemoteRockJobList) []models.ScrapedJob {
// 	var records []models.ScrapedJob
// 	var singleJob RemoteRockJob
// 	props := responseJSON.Props
// 	pageProps := props.PageProps
// 	jobs := pageProps.InitialJobOpenings

// 	for jobID, job := range jobs {

// 		rrs.Logger.Printf("[%s] Parsing job posting %d/%d", rrs.ScraperName, jobID+1, len(jobs))

// 		companyInfo := job.Company
// 		companyName := strings.ReplaceAll(companyInfo.Name, " ", "-")
// 		company := utils.NormalizeString(companyName)
// 		companyImageURL := companyInfo.ProfilePicURL
// 		slug := job.Slug

// 		jobURL := fmt.Sprintf("https://www.remoterocketship.com/company/%s/jobs/%s/", companyName, slug)

// 		if utils.IsForbidden(company, constants.ForbiddenCompanies) {
// 			rrs.Logger.Printf("Discarding company %s", company)
// 			continue
// 		}

// 		resp, err := rrs.jobInformationRequest(jobURL)
// 		if err != nil {
// 			log.Printf("error fetching job info: %w", err)
// 			continue
// 		}
// 		defer resp.Body.Close()

// 		bodyBytes, err := io.ReadAll(resp.Body)
// 		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
// 		if err != nil {
// 			log.Fatal("error parsing job HTML: %w", err)
// 		}

// 		utils.GetJSONFromHTML(doc, &singleJob)

// 		p := singleJob.Props
// 		jobInfoProps := p.PageProps

// 		title := job.RoleTitle
// 		jobOpening := jobInfoProps.JobOpening
// 		description := strings.Join([]string{
// 			jobOpening.RoleDescription,
// 			jobOpening.RoleRequirements,
// 			jobOpening.Benefits,
// 		}, "\n")

// 		location := job.Location
// 		datePosted := job.CreatedAt

// 		keyword := utils.FindKeywordInDescription(description)

// 		record := models.ScrapedJob{
// 			Title:           title,
// 			Company:         company,
// 			Location:        location,
// 			URL:             job.URL,
// 			DatePosted:      utils.FromTimestampToISOFormat(datePosted.UnixMilli()),
// 			Modality:        "Remote",
// 			Platform:        rrs.ScraperName,
// 			KeywordAppeared: keyword,
// 			LogoURL:         companyImageURL,
// 		}
// 		records = append(records, record)

// 	}

// 	return records
// }

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
