package infojobs

import "time"

type IjData struct {
	Offers []struct {
		Code         string    `json:"code"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		City         string    `json:"city"`
		Link         string    `json:"link"`
		ContractType string    `json:"contractType"`
		Workday      string    `json:"workday"`
		Teleworking  string    `json:"teleworking"`
		PublishedAt  time.Time `json:"publishedAt"`
		CompanyName  string    `json:"companyName"`
		CompanyLogo  string    `json:"companyLogo"`
		CompanyLink  string    `json:"companyLink"`
		States       []any     `json:"states"`
		Upsellings   []any     `json:"upsellings"`
		Executive    bool      `json:"executive"`
		NewBOID      string    `json:"newBOId"`
		Salary       struct {
			Range struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"range"`
			Period   string `json:"period"`
			Currency string `json:"currency"`
			Type     string `json:"type"`
		} `json:"salary,omitempty"`
	} `json:"offers"`
	Aggregation struct {
		Province []struct {
			Value        string `json:"value"`
			Label        string `json:"label"`
			Count        int    `json:"count"`
			Order        int    `json:"order"`
			SemanticLink string `json:"semanticLink"`
		} `json:"province"`
		City []struct {
			Value string `json:"value"`
			Label string `json:"label"`
			Count int    `json:"count"`
			Order int    `json:"order"`
		} `json:"city"`
		Teleworking []struct {
			Value string `json:"value"`
			Label string `json:"label"`
			Count int    `json:"count"`
			Order int    `json:"order"`
		} `json:"teleworking"`
		Category []struct {
			Value        string `json:"value"`
			Label        string `json:"label"`
			Count        int    `json:"count"`
			Order        int    `json:"order"`
			SemanticLink string `json:"semanticLink"`
		} `json:"category"`
		Education []struct {
			Value string `json:"value"`
			Label string `json:"label"`
			Count int    `json:"count"`
			Order int    `json:"order"`
		} `json:"education"`
		Workday []struct {
			Value string `json:"value"`
			Label string `json:"label"`
			Count int    `json:"count"`
			Order int    `json:"order"`
		} `json:"workday"`
		ContractType []struct {
			Value string `json:"value"`
			Label string `json:"label"`
			Count int    `json:"count"`
			Order int    `json:"order"`
		} `json:"contractType"`
		Segment []any `json:"segment"`
		Country []struct {
			Value string `json:"value"`
			Label string `json:"label"`
			Count int    `json:"count"`
			Order int    `json:"order"`
		} `json:"country"`
	} `json:"aggregation"`
	Navigation struct {
		Self                    int      `json:"self"`
		Pages                   []int    `json:"pages"`
		TotalPages              int      `json:"totalPages"`
		TotalElements           int      `json:"totalElements"`
		SortedBy                string   `json:"sortedBy"`
		AvailableSortingMethods []string `json:"availableSortingMethods"`
	} `json:"navigation"`
	Search struct {
		Keyword struct {
			Value                 string `json:"value"`
			NormalizedJobTitleIds []any  `json:"normalizedJobTitleIds"`
		} `json:"keyword"`
		CountryIds         []any  `json:"countryIds"`
		ProvinceIds        []any  `json:"provinceIds"`
		CityIds            []any  `json:"cityIds"`
		TeleworkingIds     []any  `json:"teleworkingIds"`
		CategoryIds        []any  `json:"categoryIds"`
		SubcategoryIds     []any  `json:"subcategoryIds"`
		EducationIds       []any  `json:"educationIds"`
		WorkdayIds         []any  `json:"workdayIds"`
		ContractTypeIds    []any  `json:"contractTypeIds"`
		SegmentID          any    `json:"segmentId"`
		Salary             any    `json:"salary"`
		Experience         any    `json:"experience"`
		SinceDate          string `json:"sinceDate"`
		OnlyForeignCountry bool   `json:"onlyForeignCountry"`
		SortBy             string `json:"sortBy"`
		Location           struct {
			Countries []any `json:"countries"`
			Provinces []any `json:"provinces"`
		} `json:"location"`
	} `json:"search"`
	Overview struct {
		TotalElements int    `json:"totalElements"`
		Title         string `json:"title"`
		Keyword       struct {
			Value string `json:"value"`
			Link  string `json:"link"`
			Title string `json:"title"`
		} `json:"keyword"`
		Countries []any `json:"countries"`
		Provinces []any `json:"provinces"`
		Cities    []any `json:"cities"`
	} `json:"overview"`
	Metadata struct {
		Title             string `json:"title"`
		MetaKeywords      string `json:"metaKeywords"`
		MetaDescription   string `json:"metaDescription"`
		CanonicalURL      string `json:"canonicalUrl"`
		ActiveExperiments []any  `json:"activeExperiments"`
	} `json:"metadata"`
	Segmentation struct {
		UtagParams struct {
			EsSchAdProvince          []string `json:"es_sch_ad_province"`
			EsSchAdProvinceID        []int    `json:"es_sch_ad_province_id"`
			EsSchFilterAutocomplete  int      `json:"es_sch_filter_autocomplete"`
			EsSchSearchResultsNumber int      `json:"es_sch_search_results_number"`
			EsSchSearchTerms         string   `json:"es_sch_search_terms"`
			OfferSearchPage          int      `json:"offer_search_page"`
			OfferSearchPageLimit     int      `json:"offer_search_page_limit"`
			RegionLevel2ID           []int    `json:"region_level2_id"`
			SearchTerms              string   `json:"search_terms"`
			AdsNamePage              string   `json:"ads_name_page"`
			Seccion                  string   `json:"seccion"`
			Site                     string   `json:"site"`
			IsMsite                  string   `json:"is_msite"`
			EsSchAdsNamePage         string   `json:"es_sch_ads_name_page"`
			EsSchEventName           string   `json:"es_sch_event_name"`
			EsSchAdobeIsEnabled      string   `json:"es_sch_adobe_is_enabled"`
		} `json:"utagParams"`
	} `json:"segmentation"`
	CompanyLogos struct {
		Referer      string `json:"referer"`
		Segmentation struct {
			Keyword []string `json:"keyword"`
		} `json:"segmentation"`
	} `json:"companyLogos"`
	EligibleForAutomaticAlertCreation bool `json:"eligibleForAutomaticAlertCreation"`
	TrackingProperties                struct {
		PageType                   string `json:"page_type"`
		Channel                    string `json:"channel"`
		OfferOrigin                string `json:"offer_origin"`
		Vertical                   string `json:"vertical"`
		Section                    string `json:"section"`
		Platform                   string `json:"platform"`
		CategoryLevel1             string `json:"category_level1"`
		SearchWithAutocomplete     bool   `json:"search_with_autocomplete"`
		Executive                  bool   `json:"executive"`
		Site                       string `json:"site"`
		Province                   string `json:"province"`
		ListResultNum              int    `json:"list_result_num"`
		PageName                   string `json:"page_name"`
		EventName                  string `json:"event_name"`
		Keyword                    string `json:"keyword"`
		International              bool   `json:"international"`
		SegmentEventAlreadyTracked bool   `json:"segment_event_already_tracked"`
		Order                      string `json:"order"`
	} `json:"trackingProperties"`
}
