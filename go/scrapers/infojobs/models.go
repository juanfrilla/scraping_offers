package infojobs

type Root struct {
	User            interface{}        `json:"user"`
	Experimentation Experimentation    `json:"experimentation"`
	HeaderLinks     HeaderLinks        `json:"headerLinks"`
	FooterLinks     FooterLinks        `json:"footerLinks"`
	Notification    NotificationParams `json:"notificationParams"`
	Offers          []Offer            `json:"offers"`
	Aggregation     Aggregation        `json:"aggregation"`
	Navigation      Navigation         `json:"navigation"`
	Search          Search             `json:"search"`
	Overview        Overview           `json:"overview"`
	Metadata        Metadata           `json:"metadata"`
	Segmentation    Segmentation       `json:"segmentation"`
	CompanyLogos    CompanyLogos       `json:"companyLogos"`
	TrackingProps   TrackingProperties `json:"trackingProperties"`
	TopSearches     []TopSearch        `json:"topSearches"`
	TrackingTags    map[string]Tag     `json:"trackingTags"`
	Device          string             `json:"device"`
	SearchByType    string             `json:"searchByType"`
}

type Experimentation struct {
	AnonymousID   string `json:"anonymousId"`
	SegmentUserID string `json:"segmentUserId"`
	CurrentUserID string `json:"currentUserId"`
}

type HeaderLinks struct {
	Navigation []Link `json:"navigation"`
	User       []Link `json:"user"`
}

type FooterLinks struct {
	Portals    []Link      `json:"portals"`
	Social     SocialLinks `json:"social"`
	Navigation FooterNav   `json:"navigation"`
	AppsStore  AppsStore   `json:"appsStore"`
}

type SocialLinks struct {
	Facebook Link `json:"facebook"`
	Twitter  Link `json:"twitter"`
	Youtube  Link `json:"youtube"`
}

type FooterNav struct {
	Us    []Link `json:"us"`
	About []Link `json:"about"`
	More  []Link `json:"more"`
	Press []Link `json:"press"`
}

type AppsStore struct {
	Android AppStoreItem `json:"android"`
	IOS     AppStoreItem `json:"ios"`
}

type AppStoreItem struct {
	Href string `json:"href"`
	Img  string `json:"img"`
	Alt  string `json:"alt"`
}

type NotificationParams struct {
	OfferName string `json:"offerName"`
}

type Offer struct {
	Code         string   `json:"code"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	City         string   `json:"city"`
	Link         string   `json:"link"`
	ContractType string   `json:"contractType"`
	Salary       *Salary  `json:"salary"`
	Workday      string   `json:"workday"`
	Teleworking  string   `json:"teleworking"`
	PublishedAt  string   `json:"publishedAt"`
	CompanyName  string   `json:"companyName"`
	CompanyLogo  string   `json:"companyLogo"`
	CompanyLink  string   `json:"companyLink"`
	States       []string `json:"states"`
	Upsellings   []string `json:"upsellings"`
	Executive    bool     `json:"executive"`
	NewBOId      string   `json:"newBOId"`
}

type Salary struct {
	Range struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"range"`
	Period   string `json:"period"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
}

type Aggregation struct {
	Province     []AggItem `json:"province"`
	City         []AggItem `json:"city"`
	Teleworking  []AggItem `json:"teleworking"`
	Category     []AggItem `json:"category"`
	Education    []AggItem `json:"education"`
	Workday      []AggItem `json:"workday"`
	ContractType []AggItem `json:"contractType"`
	Segment      []AggItem `json:"segment"`
	Country      []AggItem `json:"country"`
}

type AggItem struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	Count        int    `json:"count"`
	Order        int    `json:"order"`
	SemanticLink string `json:"semanticLink,omitempty"`
}

type Navigation struct {
	Self             int      `json:"self"`
	Pages            []int    `json:"pages"`
	TotalPages       int      `json:"totalPages"`
	TotalElements    int      `json:"totalElements"`
	SortedBy         string   `json:"sortedBy"`
	AvailableSorting []string `json:"availableSortingMethods"`
}

type Search struct {
	Keyword struct {
		Value string `json:"value"`
	} `json:"keyword"`
	CountryIds     []string `json:"countryIds"`
	ProvinceIds    []string `json:"provinceIds"`
	CityIds        []string `json:"cityIds"`
	TeleworkingIds []string `json:"teleworkingIds"`
}

type Overview struct {
	TotalElements int `json:"totalElements"`
}

type Metadata struct {
	Title             string   `json:"title"`
	MetaKeywords      string   `json:"metaKeywords"`
	MetaDescription   string   `json:"metaDescription"`
	CanonicalURL      string   `json:"canonicalUrl"`
	ActiveExperiments []string `json:"activeExperiments"`
}

type Segmentation struct {
	UtagParams map[string]interface{} `json:"utagParams"`
}

type CompanyLogos struct {
	Referer      string              `json:"referer"`
	Segmentation map[string][]string `json:"segmentation"`
}

type TrackingProperties struct {
	Country  string `json:"country"`
	PageType string `json:"page_type"`
	Keyword  string `json:"keyword"`
}

type TopSearch struct {
	Title string `json:"title"`
	Links []Link `json:"links"`
}

type Tag struct {
	Track           string `json:"track"`
	TrackProperties string `json:"track-properties"`
}

type Link struct {
	Href      string `json:"href"`
	Title     string `json:"title"`
	Label     string `json:"label,omitempty"`
	ClassName string `json:"className,omitempty"`
	Type      string `json:"type,omitempty"`
}
