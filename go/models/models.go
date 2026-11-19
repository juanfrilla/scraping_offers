package models

import "time"

type ScrapedJob struct {
	Title           string    `json:"title"`
	Company         string    `json:"company"`
	Location        string    `json:"location"`
	URL             string    `json:"url"`
	DatePosted      time.Time `json:"date_posted"`
	Modality        string    `json:"modality"`
	Platform        string    `json:"platform"`
	KeywordAppeared string    `json:"keyword_appeared"`
	LogoURL         string    `json:"logo_url"`
	ExternalID      string    `json:"external_id"`
}

type SavedJob struct {
	ID int64 `json:"id"`
	ScrapedJob
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



// aqui lo de is active no, quizas las updatedAt que no se actualizaron hoy se eliminan