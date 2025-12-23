import html
import json
import random
import time

import curl_cffi as requests
from bs4 import BeautifulSoup

from logger import get_logger
from utils import (
    determine_modality,
    keyword_counter,
    load_html,
    normalize_string,
)


class LinkedinScraper:
    def __init__(self):
        self.logger = get_logger("linkedin_scraper")
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

    def retrieve_offers(self) -> BeautifulSoup:
        impersonate = self.get_impersonator()
        self.logger.info(f"using {impersonate}")
        self.impersonate = impersonate
        response = self.linkedin_jobsearch_request()
        return BeautifulSoup(response.text, "html.parser")

    def get_impersonator(self):
        impersonate = random.choice(self.IMPERSONATE_LIST)
        self.logger.info(
            f"Selected LinkedIn impersonate for scraping is {impersonate}."
        )
        return impersonate

    def make_request(self, url: str, max_retries: int = 5):
        retries = 0

        while retries < max_retries:
            try:
                response = self.session.get(
                    url,
                    impersonate=self.impersonate,
                )

                if response.status_code == 429:
                    # Too many requests, wait and regenerate profile
                    wait = random.randint(1, 5)
                    self.logger.info(
                        f"Rate limited. Generating new profile, waiting {wait} seconds..."
                    )
                    time.sleep(wait)
                    self.impersonate = self.get_impersonator()
                    self.session.cookies = {}
                    self.linkedin_jobsearch_request()
                    retries += 1
                    continue  # retry the request with new profile

                # Success or other status code
                return response

            except requests.exceptions.RequestException as e:
                # Handle network errors
                wait = random.randint(1, 5)
                self.logger.warning(
                    f"Request failed ({e}), retrying in {wait} seconds..."
                )
                time.sleep(wait)
                retries += 1

        # Max retries exceeded
        self.logger.error(f"Failed to fetch {url} after {max_retries} retries.")
        return None

    def linkedin_jobsearch_request(self):
        url = "https://www.linkedin.com/jobs/search/?geoId=105646813&keywords=%22scraping%22%20OR%20%22crawling%22%20OR%20%22data%20extraction%22%20OR%20%22data%20acquisition%22%20OR%20%22data%20extraction%22&origin=JOB_SEARCH_PAGE_JOB_FILTER&refresh=true&sortBy=R"
        return self.make_request(url)

    def parse(self, soup: BeautifulSoup) -> list:
        job_list = soup.select("ul.jobs-search__results-list > li")
        self.logger.info(f"Found {len(job_list)} job postings on the page.")

        records = []

        for job_id, job in enumerate(job_list):
            url = job.select_one("a")["href"].strip()
            self.logger.info(f"Parsing job posting {job_id + 1}/{len(job_list)}")
            offer_response = self.make_request(url)
            offer_soup = BeautifulSoup(offer_response.text, "html.parser")
            script_tag = offer_soup.find("script", type="application/ld+json")
            if script_tag:
                raw_json_content = script_tag.string
                json_content = self.json_from_ld(raw_json_content)
                title = json_content.get("title", "N/A")
                company = json_content.get("hiringOrganization", {}).get("name", "N/A")
                location = (
                    json_content.get("jobLocation", {})
                    .get("address", {})
                    .get("addressLocality", "N/A")
                )
                date_posted_str = json_content.get("datePosted", "")
                description = json_content.get("description", "")
                modality = determine_modality(title, description)
            else:
                raw_title = offer_soup.find("h1", class_="topcard__title")
                title = self.get_text(raw_title)
                raw_company = offer_soup.find("a", class_="topcard__org-name-link")
                company = self.get_text(raw_company)
                raw_location = offer_soup.find("span", class_="topcard__flavor--bullet")
                location = self.get_text(raw_location)
                raw_date_posted_str = offer_soup.find("time")["datetime"].strip()
                date_posted_str = f"{raw_date_posted_str}T00:00:00Z"
                raw_description = offer_soup.find("div", class_="description__text")
                description = self.get_text(raw_description)

                modality = determine_modality(title, description)
            records.append(
                {
                    "title": title,
                    "company": normalize_string(company),
                    "location": self.determine_location(location),
                    "url": url,
                    "date_posted": date_posted_str,
                    "modality": modality,
                    "platform": "LINKEDIN",
                    "keywords": keyword_counter(description),
                }
            )

        return records

    def get_text(self, soup: BeautifulSoup) -> str:
        if soup:
            return soup.get_text().strip()
        return ""

    def determine_location(self, location: str) -> str:
        locations = {
            "Sevilla La Nueva": "Sevilla",
            "Comunidad De Madrid, España": "Madrid",
            "Valencia/València": "Valencia",
            "Alcobendas": "Madrid",
        }
        location_lower = location.lower()
        if location_lower.endswith("y alrededores"):
            location_lower = location_lower.replace("y alrededores", "")
        normalized = normalize_string(location_lower)
        return locations.get(normalized, normalized)

    def json_from_ld(self, raw: str) -> dict:
        clean = raw.strip()
        clean = html.unescape(clean)
        return json.loads(clean)

    def scrape(self):
        self.logger.info("Starting Linkedin scraping.")
        soup = self.retrieve_offers()
        jobs = self.parse(soup)
        self.logger.info("Finished Linkedin scraping.")
        return jobs

    def scrape_test(self):
        html_content = load_html("./seed/linkedin_jobsearch_more_offers.html")
        return self.parse(html_content)
