package remoteok

import "encoding/json"

type RemoteOkJob struct {
	Title          string   `json:"title"`
	DatePosted     string   `json:"datePosted"`
	Description    string   `json:"description"`
	Image          string   `json:"image"`
	EmploymentType string   `json:"employmentType"`
	JobBenefits    []string `json:"jobBenefits"`
	Industry       string   `json:"industry"`
	ValidThrough   string   `json:"validThrough"`

	BaseSalary struct {
		Currency string `json:"currency"`
		Value    struct {
			MinValue json.Number `json:"minValue"`
			MaxValue json.Number `json:"maxValue"`
			UnitText string      `json:"unitText"`
		} `json:"value"`
	} `json:"baseSalary"`

	HiringOrg struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"hiringOrganization"`

	LocationReqs []struct {
		Name string `json:"name"`
	} `json:"applicantLocationRequirements"`
}

