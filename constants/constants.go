package constants

const DataFile = "./data/jobs.json"

var ForbiddenKeywords = []string{
	"account",
	"ai",
	"artificial intelligence",
	"bi",
	"brand",
	"business development",
	"business intelligence",
	"content",
	"deal",
	"director",
	"editor",
	"entry level",
	"executive",
	"field",
	"founder",
	"head",
	"intern",
	"laborer",
	"lead",
	"manager",
	"marketing",
	"operation",
	"principal",
	"product",
	"recruiter",
	"sale",
	"sap",
	"scientist",
	"search",
	"seo",
	"strateg",
	"talent aquisition",
	"technician",
	"technisian",
	"veterinarian",
}

// SearchKeywords - Terms used to query the Empleate API.
var SearchKeywords = []string{
	"scraping",
	"crawling",
	"data%20aquisition",
	"scraper",
	"crawler",
	"anti-bot",
	"antibot",
}

// ForbiddenCompanies - Companies whose job postings should be filtered out.
var ForbiddenCompanies = []string{
	"infatica.io",
	"corsearch",
	"mindrift",
	"wayops",
	"duckduckgo",
	"fundraise up",
}

var JsonTypes = []string{"application/json", "application/ld+json"} // constants

var ImpersonateList = []string{
	// Chrome Desktop
	"chrome-133",
	"chrome-141",
	"chrome-143",
	"chrome-144",
	"chrome-145",
	"chrome-146",
	"chrome-147",
	"chrome-148",

	// Chrome per-OS
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

	// Chrome latest aliases
	"chrome-latest",
	"chrome-latest-windows",
	"chrome-latest-linux",
	"chrome-latest-macos",
	"chrome-latest-android",
	"chrome-latest-ios",

	// Firefox
	"firefox-133",
	"firefox-148",

	// // Safari Desktop
	// "safari-18",
	// "safari-17",
	// "safari-16",

	// // Safari iOS
	// "safari-17-ios",
	// "safari-18-ios",
}
