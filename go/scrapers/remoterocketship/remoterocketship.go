package remoterocketship

import (
	"fmt"
	"io"
	"log"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type RemoteRocketshipScraper struct {
	Session     tls_client.HttpClient
	Logger      *log.Logger
	ScraperName string
}

func NewRemoteRocketshipScraper() *RemoteRocketshipScraper {
	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_146),
		tls_client.WithInsecureSkipVerify(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	scraperName := "RemoteRocketship"
	return &RemoteRocketshipScraper{
		Session:     client,
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
	}
}

func (rrs *RemoteRocketshipScraper) getJobsRequest() (*http.Response, error) {
	url := "https://www.remoterocketship.com/?page=1&sort=DateAdded&jobTitle=scraping%2Ccrawling%2Cscraper%2Ccrawler%2Cdata+acquisition&locations=Europe%2CSpain"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return rrs.Session.Do(req)
}

func (rrs *RemoteRocketshipScraper) jobInformationRequest(jobURL string) (*http.Response, error) {

	req, err := http.NewRequest("GET", jobURL, nil)
	if err != nil {
		log.Fatal("Error creating request: %w", err)
	}

	req.Header.Set("referer", "https://www.remoterocketship.com/?page=1&sort=DateAdded&jobTitle=scraping%2Ccrawling%2Cscraper%2Ccrawler%2Cdata+acquisition&locations=Europe%2CSpain")
	return rrs.Session.Do(req)

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
