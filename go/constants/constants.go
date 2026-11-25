package constants

import (
	"fmt"
	"os"
	"path/filepath"
)

var ForbiddenKeywords = []string{
	"account",
	"AEO/GEO",
	"ai",
	"architect",
	"artificial intelligence",
	"automotive",
	"bi",
	"brand",
	"business analyst",
	"business development",
	"business intelligence",
	"chemical",
	"commercial",
	"compliance",
	"content",
	"corporate",
	"customer support",
	"deal",
	"devops",
	"director",
	"economist",
	"editor",
	"electrical",
	"electronics",
	"embedded",
	"entry level",
	"executive",
	"field",
	"finance",
	"financial",
	"founder",
	"front end",
	"front-end",
	"frontend",
	// "full stack",
	// "full-stack",
	// "fullstack",
	"graphics",
	"growth",
	"gtm",
	"hardware",
	"head",
	"human resources",
	"hydraulics",
	"intern",
	"laborer",
	"lead",
	"machine learning",
	"ml",
	"manager",
	"market",
	"marketing",
	"mecatronics",
	"mechanical",
	"nuclear",
	"operation",
	"operator",
	"plasma",
	"principal",
	"product",
	"propulsion",
	"qa",
	"quality assurance",
	"recruiter",
	"rf",
	"robotics",
	"rocket",
	"sale",
	"sales",
	"sap",
	"science",
	"scientist",
	"search",
	"seo",
	"simulation",
	"solutions",
	"strateg",
	"supply",
	"system",
	"talent",
	"tax",
	"technician",
	"technisian",
	"test",
	"thermal",
	"veterinarian",
	//"analytics",
	//"analyst",
}

var SearchKeywords = []string{
	"scraping",
	"crawling",
	"data%20aquisition",
	"scraper",
	"crawler",
	"anti-bot",
	"antibot",
	"anti%20bot",
	"data%20extraction",
	"reverse%20engineering",
	"security20%researcher",
}

var ForbiddenCompanies = []string{
	"infatica.io",
	"mindrift",
	"duckduckgo",
	"searchapi",
	"apify",
}

var ForbiddenCountries = []string{
	"india", "philippines", "brazil", "argentina", "mexico", "panama", "korea",
}

var JsonTypes = []string{"application/json", "application/ld+json"}

var ImpersonateList = []string{
	"chrome-133",
	"chrome-141",
	"chrome-143",
	"chrome-144",
	"chrome-145",
	"chrome-146",
	"chrome-147",
	"chrome-148",

	"chrome-133-windows",
	"chrome-133-linux",
	"chrome-133-macos",
	"chrome-133-android",
	"chrome-133-ios",

	"chrome-146-windows",
	"chrome-146-linux",
	"chrome-146-macos",
	"chrome-146-android",
	"chrome-146-ios",

	"chrome-latest",
	"chrome-latest-windows",
	"chrome-latest-linux",
	"chrome-latest-macos",
	"chrome-latest-android",
	"chrome-latest-ios",

	"firefox-133",
	"firefox-148",
}

var (
	BaseDir         string
	LogDir          string
	DataFile        string
	DataDB          string
	LastScrapedFile string
	EnvFile         string
)

func InitPaths() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	currentDir, err := filepath.Abs(workingDir)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			BaseDir = currentDir
			break
		}

		parent := filepath.Dir(currentDir)

		if parent == currentDir {
			return fmt.Errorf("could not find go.mod starting from %s", workingDir)
		}

		currentDir = parent
	}

	LogDir = filepath.Join(BaseDir, "logs")
	DataFile = filepath.Join(BaseDir, "data", "jobs.json")
	DataDB = filepath.Join(BaseDir, "data", "jobs.db")
	LastScrapedFile = filepath.Join(BaseDir, "data", "last_scraped.json")
	EnvFile = filepath.Join(BaseDir, ".env")

	return nil
}
