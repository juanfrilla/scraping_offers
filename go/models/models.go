package models

type ScrapedJob struct {
	Title           string `json:"title"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	URL             string `json:"url"`
	DatePosted      string `json:"date_posted"`
	Modality        string `json:"modality"`
	Platform        string `json:"platform"`
	KeywordAppeared string `json:"keyword_appeared"`
	LogoURL         string `json:"logo_url"`
}
