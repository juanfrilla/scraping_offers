import random
import time

import curl_cffi as requests
from bs4 import BeautifulSoup

from logger import get_logger
from utils import (
    determine_modality,
    filter_jobs,
    find_keyword,
    from_timestamp_to_isoformat,
    get_json_from_html,
    normalize_string,
)


class SimplyHiredScraper:
    # Hacen falta residential proxies
    def __init__(self):
        self.session = requests.Session()
        self.impersonate = None

        self.IMPERSONATE_LIST = [
            # Chrome Desktop
            "chrome124",
            "chrome131",
            "chrome133a",
            "chrome136",
            # Chrome Android
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
        self.scraper_name = "SimplyHired"
        self.logger = get_logger(self.scraper_name)

    def get_jobs_request(self, keyword: str):
        burp0_url = f"https://www.simplyhired.es/search?q={keyword}&l=espa%25c3%25b1a"
        for impersonate in self.IMPERSONATE_LIST:
            response = self.session.get(
                "https://www.simplyhired.es/", impersonate=impersonate
            )
            self.impersonate = impersonate
            self.logger.info(f"Using {impersonate}")
            response = self.session.get(burp0_url, impersonate=impersonate)
            html_json = get_json_from_html(BeautifulSoup(response.text, "html.parser"))
            if html_json != {}:
                self.logger.info("Retrieved data!")
                return html_json
            wait = random.randint(1, 5)
            self.logger.info(f"Empty data, continuing and waiting {wait} seconds ....")
            time.sleep(wait)
        return {}

    def job_info_request(self, bot_url: str):
        url = f"https://www.simplyhired.es{bot_url}"
        response = self.session.get(url, impersonate=self.impersonate)
        offer_soup = BeautifulSoup(response.text, "html.parser")
        return get_json_from_html(offer_soup)

    def parse(self, jobs_data: dict) -> list:
        records = []
        jobs = jobs_data.get("props", {}).get("pageProps", {}).get("jobs", [])

        for job_id, job in enumerate(jobs):
            self.logger.info(f"Parsing job posting {job_id + 1}/{len(jobs)}")
            bot_url = job.get("botUrl")
            job_info = self.job_info_request(bot_url)
            job_info_props = job_info.get("props", {}).get("pageProps", {})
            title = job.get("title", "N/A")
            company = job.get("company")

            description = (
                job_info.get("description")
                or job_info_props.get("jobDescriptionHtml")
                or "N/A"
            )

            location = job.get("location")
            date_posted_timestamp = job.get("datePublished") or job.get("dateOnIndeed")
            modality = determine_modality(title, description)
            keyword_appeared = find_keyword(description, title)
            if keyword_appeared:
                records.append(
                    {
                        "title": title,
                        "company": normalize_string(company),
                        "location": location,
                        "url": f"https://www.simplyhired.es{bot_url}",
                        "date_posted": from_timestamp_to_isoformat(
                            date_posted_timestamp
                        ),
                        "modality": modality,
                        "platform": "SIMPLYHIRED",
                        "keyword_appeared": keyword_appeared,
                        "logo_url": job_info_props.get("employerSquareLogoUrl", ""),
                    }
                )

        return records

    def scrape(self):
        keywords = [
            "scraping",
            "crawling",
            "data%20aquisition",
            "data%20extraction",
            "scraper",
            "crawler",
        ]
        all_jobs = []
        for keyword in keywords:
            self.logger.info(f"Scraping word {keyword}")
            jobs_data = self.get_jobs_request(keyword)
            jobs = self.parse(jobs_data)
            self.logger.info(f"Retrieved {len(jobs)} for {keyword}")
            all_jobs += jobs
        return filter_jobs(all_jobs)
