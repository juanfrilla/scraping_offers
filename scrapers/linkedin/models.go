package linkedin

import "time"

type LinkedInJob struct {
	Context            string    `json:"@context"`
	Type               string    `json:"@type"`
	DatePosted         time.Time `json:"datePosted"`
	Description        string    `json:"description"`
	EmploymentType     string    `json:"employmentType"`
	HiringOrganization struct {
		Type   string `json:"@type"`
		Name   string `json:"name"`
		SameAs string `json:"sameAs"`
		Logo   string `json:"logo"`
	} `json:"hiringOrganization"`
	Identifier struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"identifier"`
	Image       string `json:"image"`
	Industry    string `json:"industry"`
	JobLocation struct {
		Type    string `json:"@type"`
		Address struct {
			Type            string `json:"@type"`
			AddressCountry  string `json:"addressCountry"`
			AddressLocality string `json:"addressLocality"`
			AddressRegion   any    `json:"addressRegion"`
			StreetAddress   any    `json:"streetAddress"`
		} `json:"address"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"jobLocation"`
	Skills                string    `json:"skills"`
	Title                 string    `json:"title"`
	ValidThrough          time.Time `json:"validThrough"`
	EducationRequirements struct {
		Type               string `json:"@type"`
		CredentialCategory string `json:"credentialCategory"`
	} `json:"educationRequirements"`
	ExperienceRequirements struct {
		Type               string `json:"@type"`
		MonthsOfExperience int    `json:"monthsOfExperience"`
	} `json:"experienceRequirements"`
	JobLocationType               string `json:"jobLocationType"`
	ApplicantLocationRequirements struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"applicantLocationRequirements"`
}
