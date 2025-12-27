import curl_cffi as requests
from bs4 import BeautifulSoup

from logger import get_logger
from utils import (
    determine_modality,
    from_timestamp_to_isoformat,
    get_json_from_html,
    keyword_counter,
    normalize_string,
)


class SimplyHiredScraper:
    def __init__(self):
        self.logger = get_logger("simplyhired_scraper")
        self.session = requests.Session()

        self.IMPERSONATE_LIST = [
            # Chrome Desktop
            "chrome99",
            "chrome100",
            "chrome101",
            "chrome104",
            "chrome107",
            "chrome110",
            "chrome116",
            "chrome119",
            "chrome120",
            "chrome123",
            "chrome124",
            "chrome131",
            "chrome133a",
            "chrome136",
            # Chrome Android
            "chrome99_android",
            "chrome131_android",
            # Edge
            "edge99",
            "edge101",
            # Safari Desktop
            "safari153",
            "safari155",
            "safari170",
            "safari180",
            "safari184",
            "safari260",
            # Safari iOS
            "safari172_ios",
            "safari180_ios",
            "safari184_ios",
            "safari260_ios",
            # Firefox
            "firefox133",
            "firefox135",
            # Tor
            "tor145",
        ]

    def get_jobs_request(self):
        burp0_url = "https://www.simplyhired.es/search?q=%22crawler%22+or+%22scraping%22&l=espa%25c3%25b1a"

        response = self.session.get(burp0_url, impersonate="chrome131")
        html_json = get_json_from_html(BeautifulSoup(response.text))
        return html_json

    def job_info_request(self, bot_url: str):
        burp0_url = f"https://www.simplyhired.es{bot_url}"
        response = self.session.get(burp0_url, impersonate="chrome131")
        offer_soup = BeautifulSoup(response.text, "html.parser")
        return get_json_from_html(offer_soup)

    def parse(self, jobs_data: dict) -> list:
        records = []
        jobs = jobs_data = jobs_data.get("pageProps", {}).get("jobs", {})
        self.logger.info(f"Found {len(jobs)} job postings on the page.")

        for job_id, job in enumerate(jobs):
            self.logger.info(f"Parsing job posting {job_id + 1}/{len(jobs)}")
            bot_url = job.get("botUrl")
            job_info = self.job_info_request(bot_url)
            job_info_props = job_info.get("props", {}).get("pageProps", {})
            title = job.get("title")
            company = job.get("company")

            description = job_info.get("description") or job_info_props.get(
                "jobDescriptionHtml"
            )

            location = job.get("location")
            date_posted_timestamp = job.get("datePublished") or job.get("dateOnIndeed")
            modality = determine_modality(title, description)
            records.append(
                {
                    "title": title,
                    "company": normalize_string(company),
                    "location": location,
                    "url": bot_url,
                    "date_posted": from_timestamp_to_isoformat(date_posted_timestamp),
                    "modality": modality,
                    "platform": "SIMPLYHIRED",
                    "keywords": keyword_counter(description),
                }
            )

        return records

    def scrape(self):
        self.logger.info("Finished SimplyHired scraping.")
        jobs_data = self.get_jobs_request()
        jobs = self.parse(jobs_data)
        self.logger.info("Finished SimplyHired scraping.")
        return jobs
