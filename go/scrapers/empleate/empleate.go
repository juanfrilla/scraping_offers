package empleate

import (
	"fmt"
	"io"
	"log"
	"os"
	"scraping_offers/go/constants"
	"scraping_offers/go/models"
	"scraping_offers/go/utils"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type EmpleateScraper struct {
	Session     tls_client.HttpClient
	ScraperName string
	Logger      *log.Logger
}

func NewEmpleateScraper() *EmpleateScraper {

	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_131),
		tls_client.WithInsecureSkipVerify(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil
	}
	scraperName := "Empleate"
	return &EmpleateScraper{
		Session:     client,
		ScraperName: scraperName,
		Logger:      log.New(os.Stdout, fmt.Sprintf("[%s] ", scraperName), log.LstdFlags),
	}
}
func (es *EmpleateScraper) EmpleateJobSearchRequest(keyword string) (*fhttp.Response, error) {
	url := fmt.Sprintf("https://www.empleate.gob.es/empleate/open/offersearch/selectBuscador?q.op=AND&rows=10&sort=score%%20desc&defType=edismax&df=titulo&facet=true&facet.field=paisF&facet.field=provinciaF&facet.field=provincia&facet.field=categoria&facet.field=categoriaF&facet.field=subcategoriaF&facet.field=subcategoria&facet.field=origen&facet.field=tipoContrato&facet.field=tipoContratoN&facet.field=noMeInteresa&facet.field=educacionF&facet.field=fechaCreacionPortal&facet.field=jornadaF&facet.field=experienciaF&facet.field=educacion&facet.field=minExperiencia&facet.field=jornada&facet.field=pais&facet.field=discapacidad&facet.field=cno&facet.field=portales&facet.field=showPortalPu&facet.field=showPortalPr&facet.mincount=1&f.topics.facet.limit=50&json.nl=map&fq=(speStateId%%3A1%%20OR%%20speStateId%%3A4)%%20AND%%20checkVisible%%3A1&fl=*%%2C%%20score&q=%s&wt=json", keyword)

	req, err := fhttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return es.Session.Do(req)
}

func (es *EmpleateScraper) Parse(jsonData EmpleateJobList) []models.ScrapedJob {
	var records []models.ScrapedJob

	response := jsonData.Response
	docs := response.Docs

	for i, item := range docs {
		es.Logger.Printf("Processing offer %d/%d", i+1, len(docs))
		title := item.Titulo
		company := item.Contacto
		description := item.Contenido
		url := item.URL
		datePosted := item.FechaCreacionPortal
		platform := item.Entitytype
		location := item.ProvinciaF
		keyword := utils.FindKeywordInDescription(description)
		if keyword != "" && !utils.IsForbidden(company, constants.ForbiddenCompanies) && !utils.IsForbidden(title, constants.ForbiddenKeywords) {
			record := models.ScrapedJob{
				Title:           title,
				Company:         company,
				Location:        location,
				URL:             url,
				DatePosted:      datePosted,
				Modality:        utils.DetermineModality(title, description),
				Platform:        platform,
				KeywordAppeared: keyword,
				LogoURL:         "",
			}
			records = append(records, record)
		} else {
			es.Logger.Printf("Discarding title %s from company %s", title, company)
			continue
		}
	}

	return records
}

func (es *EmpleateScraper) Scrape() ([]models.ScrapedJob, error) {
	//TODO aqui poner los nils devolverlos
	var allJobs []models.ScrapedJob

	for _, keyword := range constants.SearchKeywords {
		es.Logger.Printf("Scraping word: %s", keyword)
		resp, err := es.EmpleateJobSearchRequest(keyword)
		if err != nil {
			return nil, fmt.Errorf("Error requesting keyword %s: %w", keyword, err)
		}

		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading body: %v", err)
		}
		var JobSearchJson EmpleateJobList
		utils.DecodeJSONP(content, &JobSearchJson)

		jobs := es.Parse(JobSearchJson)
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
}

func getString(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}
