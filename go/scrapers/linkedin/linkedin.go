package linkedin

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sardanioss/httpcloak/client"

	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"
)

type locationCfg struct {
	Name    string
	URLName string
	GeoID   string
}
type LinkedinScraper struct {
	Session     *client.Client
	ScraperName string
	Logger      *log.Logger
	impersonate string
}

func NewLinkedinScraper() *LinkedinScraper {
	scraperName := "Linkedin"
	impersonate := utils.RandomImpersonation()
	return &LinkedinScraper{
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
		Session:     client.NewClient(impersonate),
		impersonate: impersonate,
	}
}

func (ls *LinkedinScraper) Scrape() []models.ScrapedJob {
	urls, err := ls.retrieveOffers()
	if err != nil {
		log.Fatal(err)
	}

	jobs, err := ls.parse(urls)
	if err != nil {

		ls.Logger.Printf("Error parsing urls", err)

	}
	return jobs
}

func (ls *LinkedinScraper) retrieveOffers() ([]string, error) {
	var allUrls []string

	locations := []locationCfg{
		{Name: "Spain", URLName: "Espa%C3%B1a", GeoID: "105646813"},
	}

	for _, loc := range locations {
		ls.Logger.Printf("--- Starting location: %s ---", loc.Name)

		for _, keyword := range constants.SearchKeywords {

			ls.Logger.Printf("--- Looking for keyword: %s ---", keyword)
			start := 0

			jobSearchURL := fmt.Sprintf(
				"https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?keywords=%%22%s%%22&location=%s&geoId=%s&start=%d&f_WT=2",
				url.QueryEscape(keyword),
				loc.URLName,
				loc.GeoID,
				start,
			)
			resp, err := ls.makeRequest(jobSearchURL, 5)
			if err != nil {
				ls.Logger.Printf("Error getJobsRequest para %s: %v", keyword, err)
				continue
			}
			bodyBytes, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
			if err != nil {
				continue // USA CONTINUE, NO BREAK
			}

			// html, err := doc.Html()

			// htmlFilename := fmt.Sprintf("debug_%s.html", keyword)
			// err = os.WriteFile(htmlFilename, []byte(html), 0644)

			jobList := doc.Find("li")
			if jobList.Length() == 0 {
				ls.Logger.Println("No more jobs found in this pagination.")
				continue // USA CONTINUE, NO BREAK
			}
			if err != nil {
				ls.Logger.Printf("Error parsing HTML para %s: %v", keyword, err)
				continue
			}
			jobList.Each(func(_ int, s *goquery.Selection) {
				linkTag := s.Find("a").First()
				if linkTag.Length() == 0 {
					return
				}

				title := strings.ToLower(strings.TrimSpace(linkTag.Text()))
				href, exists := linkTag.Attr("href")
				if !exists {
					return
				}
				href = strings.TrimSpace(href)

				companyTag := s.Find(".base-search-card__subtitle, h4").First()
				companyName := ""
				if companyTag.Length() > 0 {
					companyName = strings.ToLower(strings.TrimSpace(companyTag.Text()))
				}

				if !utils.IsForbidden(companyName, constants.ForbiddenCompanies) && !utils.IsForbidden(title, constants.ForbiddenKeywords) {
					allUrls = append(allUrls, href)
				}

			})
		}
	}

	seen := make(map[string]struct{})
	var unique []string
	for _, u := range allUrls {
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		unique = append(unique, u)
	}

	return unique, nil
}

func (ls *LinkedinScraper) makeRequest(urlStr string, maxRetries int) (*client.Response, error) {
	retries := 0
	ctx := context.Background()
	for retries < maxRetries {
		ls.Logger.Printf("Requesting %s with impersonate %s", urlStr, ls.impersonate)
		resp, err := ls.Session.Get(ctx, urlStr, nil)
		if err != nil {
			wait := time.Duration(rand.Intn(5)+1) * time.Second
			ls.Logger.Printf("Request failed (%v), retrying in %s...", err, wait)
			time.Sleep(wait)
			retries++
			ls.impersonate = utils.RandomImpersonation()
			ls.Session = client.NewClient(ls.impersonate)
			continue
		}
		if resp.StatusCode == 429 {
			wait := time.Duration(rand.Intn(5)+1) * time.Second
			ls.Logger.Printf("Rate limited. Generating new profile, waiting %s seconds...", wait)
			time.Sleep(wait)
			retries++
			ls.impersonate = utils.RandomImpersonation()
			ls.Session = client.NewClient(ls.impersonate)
			continue
		}

		return resp, nil
	}

	ls.Logger.Printf("Failed to fetch %s after %d retries.", urlStr, maxRetries)
	return nil, fmt.Errorf("max retries exceeded for %s", urlStr)
}

func (ls *LinkedinScraper) parse(urls []string) ([]models.ScrapedJob, error) {
	var records []models.ScrapedJob
	externalIDs := make(map[string]struct{})
	todayStr := time.Now().Format("2006-01-02")

	for i, urlStr := range urls {
		ls.Logger.Printf("Parsing job posting %d/%d", i+1, len(urls))
		ls.Logger.Printf("Parsing job %s", urlStr)

		resp, err := ls.makeRequest(urlStr, 5)
		if err != nil || resp == nil || resp.StatusCode != 200 {
			continue
		}

		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
		if err != nil {
			continue
		}

		// html, err := doc.Html()

		// htmlFilename := fmt.Sprintf("debug_%s.html", "offer")
		// err = os.WriteFile(htmlFilename, []byte(html), 0644)

		var jsonContent LinkedInJob
		if err := utils.GetJSONFromHTML(doc, &jsonContent); err == nil && jsonContent.Title != "" {

			title := jsonContent.Title
			jlRaw := jsonContent.JobLocation

			location := jlRaw.Address.AddressLocality
			datePostedStr := jsonContent.DatePosted
			description := jsonContent.Description
			logoURL := jsonContent.HiringOrganization.Logo
			company := jsonContent.HiringOrganization.Name

			modality := utils.DetermineModality(title, description)
			determinedLocation := ls.determineLocation(location)
			externalID := fmt.Sprintf("%s_%s_%s", title, company, determinedLocation)
			keywordAppeared := utils.FindKeywordInDescription(description)

			if keywordAppeared != "" {
				if _, exists := externalIDs[externalID]; !exists {
					externalIDs[externalID] = struct{}{}
					records = append(records, models.ScrapedJob{
						Title:           title,
						Company:         utils.NormalizeString(company),
						Location:        determinedLocation,
						URL:             urlStr,
						DatePosted:      utils.FromTimestampToISOFormat(datePostedStr.UnixMilli()),
						Modality:        modality,
						Platform:        ls.ScraperName,
						KeywordAppeared: keywordAppeared,
						LogoURL:         logoURL,
					})
				}
			} else {
				ls.Logger.Printf("Discarding title %s from company %s", title, company)
				continue
			}

		} else {
			rawTitle := doc.Find("h1.topcard__title").First()
			title := ls.getText(rawTitle)

			rawCompany := doc.Find("a.topcard__org-name-link").First()
			company := ls.getText(rawCompany)

			rawLocation := doc.Find("span.topcard__flavor--bullet").First()
			location := ls.getText(rawLocation)

			rawDate := doc.Find("time").First()
			datePostedStr := todayStr + "T00:00:00Z"
			if rawDate.Length() > 0 {
				if dt, ok := rawDate.Attr("datetime"); ok {
					datePostedStr = strings.TrimSpace(dt) + "T00:00:00Z"
				}
			} else {
				ls.Logger.Printf("Erroneous date or not found for %s", urlStr)
			}

			rawLogo := doc.Find("img.artdeco-entity-image").First()

			logoURL := ""
			if rawLogo.Length() > 0 {
				if val, ok := rawLogo.Attr("data-delayed-url"); ok {
					logoURL = strings.TrimSpace(val)
				} else if val, ok := rawLogo.Attr("src"); ok {
					logoURL = strings.TrimSpace(val)
				}
			}

			rawDescription := doc.Find("div.description__text").First()
			description := ls.getText(rawDescription)
			modality := utils.DetermineModality(title, description)
			determinedLocation := ls.determineLocation(location)
			externalID := fmt.Sprintf("%s_%s_%s", title, company, determinedLocation)
			keywordAppeared := utils.FindKeywordInDescription(description)
			if keywordAppeared != "" {
				if _, exists := externalIDs[externalID]; !exists {
					externalIDs[externalID] = struct{}{}
					records = append(records, models.ScrapedJob{
						Title:           title,
						Company:         utils.NormalizeString(company),
						Location:        determinedLocation,
						URL:             urlStr,
						DatePosted:      datePostedStr,
						Modality:        modality,
						Platform:        ls.ScraperName,
						KeywordAppeared: keywordAppeared,
						LogoURL:         logoURL,
					})
				}
			} else {
				ls.Logger.Printf("Discarding title %s from company %s", title, company)
			}
		}
	}

	return records, nil
}

func (ls *LinkedinScraper) getText(sel *goquery.Selection) string {
	if sel == nil || sel.Length() == 0 {
		return ""
	}
	return strings.TrimSpace(sel.Text())
}

func (ls *LinkedinScraper) determineLocation(location string) string {
	mapping := map[string]string{
		"Sevilla La Nueva":                    "Sevilla",
		"Comunidad De Madrid, España":         "Madrid",
		"Valencia/València":                   "Valencia",
		"Alcobendas":                          "Madrid",
		"Greater Barcelona Metropolitan Area": "Barcelona",
		"Greater Madrid Metropolitan Area":    "Madrid",
		"Vega Del Codorno":                    "Cuenca",
	}

	locLower := strings.ToLower(location)
	if strings.HasSuffix(locLower, "y alrededores") {
		locLower = strings.TrimSpace(strings.ReplaceAll(locLower, "y alrededores", ""))
	}
	normalized := utils.NormalizeString(locLower)

	if v, ok := mapping[normalized]; ok {
		return v
	}
	return normalized
}
