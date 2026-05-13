package remoteok

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

var searchKeywords = []string{"python", "javascript"}

type RemoteOKScraper struct {
	Session     tls_client.HttpClient
	ScraperName string
	Logger      *log.Logger
}

func NewRemoteOKScraper() *RemoteOKScraper {
	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_146),
		tls_client.WithInsecureSkipVerify(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	scraperName := "RemoteOK"
	return &RemoteOKScraper{
		Session:     client,
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
	}
}

func (ro *RemoteOKScraper) getJobsRequest(keyword string) (*http.Response, error) {
	url := fmt.Sprintf(
		"https://remoteok.com/?tags=%s&action=get_jobs&premium=0&regular=1",
		keyword,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return ro.Session.Do(req)
}

func (ro *RemoteOKScraper) parse(doc *goquery.Document) []models.ScrapedJob {
	expandText := make(map[string]string)
	doc.Find("tr.expand").Each(func(_ int, row *goquery.Selection) {
		class, _ := row.Attr("class")
		for _, cls := range strings.Fields(class) {
			if strings.HasPrefix(cls, "expand-") {
				id := strings.TrimPrefix(cls, "expand-")
				expandText[id] = row.Text()
			}
		}
	})

	var singleJob RemoteOkJob
	var records []models.ScrapedJob
	jobs := doc.Find("tr.job")
	total := jobs.Length()

	jobs.Each(func(i int, job *goquery.Selection) {
		dataID, _ := job.Attr("data-id")
		botURL, _ := job.Attr("data-url")
		ro.Logger.Printf("Parsing job posting %d/%d", i+1, total)

		err := utils.GetJSONFromHTML(job, &singleJob)
		if err != nil {
			log.Fatal(err)
		}

		title := singleJob.Title
		if title == "" {
			title = "N/A"
		}
		company := singleJob.HiringOrg.Name
		datePosted := singleJob.DatePosted
		modality := "REMOTE"

		locs := make([]string, 0, len(singleJob.LocationReqs))
		for _, l := range singleJob.LocationReqs {
			locs = append(locs, l.Name)
		}
		location := strings.Join(locs, ",")

		description := expandText[dataID]
		keyword := utils.FindKeywordInDescription(description)
		if keyword != "" && !utils.IsForbidden(company, constants.ForbiddenCompanies) && !utils.IsForbidden(title, constants.ForbiddenKeywords) {
			record := models.ScrapedJob{
				Title:           title,
				Company:         utils.NormalizeString(company),
				Location:        location,
				URL:             "https://www.remoteok.com" + botURL,
				DatePosted:      datePosted,
				Modality:        modality,
				Platform:        ro.ScraperName,
				KeywordAppeared: keyword,
				LogoURL:         singleJob.Image,
			}
			records = append(records, record)
		}
	})

	return records
}

func (s *RemoteOKScraper) Scrape() []models.ScrapedJob {
	var allJobs []models.ScrapedJob

	for _, keyword := range searchKeywords {
		s.Logger.Printf("Scraping keyword: %s", keyword)

		resp, err := s.getJobsRequest(keyword)
		if err != nil {
			s.Logger.Printf("Request error for %s: %v", keyword, err)
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			s.Logger.Printf("Read error for %s: %v", keyword, err)
			continue
		}
		wrapped := "<table>" + strings.ReplaceAll(string(bodyBytes), "\t", " ") + "</table>"

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrapped))
		if err != nil {
			s.Logger.Printf("Parse error for %s: %v", keyword, err)
			continue
		}

		jobs := s.parse(doc)
		s.Logger.Printf("Retrieved %d jobs for %s", len(jobs), keyword)
		allJobs = append(allJobs, jobs...)
	}
	return allJobs
}
