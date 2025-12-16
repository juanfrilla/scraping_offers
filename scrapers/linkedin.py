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
        self.profile = self.get_profile()

    def get_profile(self):
        profiles = [
            {
                "impersonate": "firefox135",
                "headers": {
                    "User-Agent": (
                        "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) "
                        "Gecko/20100101 Firefox/135.0"
                    ),
                    "Accept": (
                        "text/html,application/xhtml+xml,application/xml;"
                        "q=0.9,image/avif,image/webp,*/*;q=0.8"
                    ),
                    "Accept-Language": "es-ES,es;q=0.9",
                    "Accept-Encoding": "gzip, deflate, br, zstd",
                    "Upgrade-Insecure-Requests": "1",
                    "Sec-Fetch-Site": "none",
                    "Sec-Fetch-Mode": "navigate",
                    "Sec-Fetch-User": "?1",
                    "Sec-Fetch-Dest": "document",
                },
            },
        ]

        profile = random.choice(profiles)
        self.logger.info(
            f"Selected LinkedIn profile for scraping is {profile.get('impersonate')}."
        )
        return profile

    def make_request(self, url: str, max_retries: int = 5):
        retries = 0

        while retries < max_retries:
            try:
                response = self.session.get(
                    url,
                    headers=self.profile["headers"],
                    impersonate=self.profile.get("impersonate"),
                )

                if response.status_code == 429:
                    # Too many requests, wait and regenerate profile
                    wait = random.randint(1, 5)
                    self.logger.info(
                        f"Rate limited. Generating new profile, waiting {wait} seconds..."
                    )
                    time.sleep(wait)
                    self.profile = self.get_profile()
                    self.session.cookies = {}
                    self.linkedin_jobsearch_request()
                    retries += 1
                    continue  # retry the request with new profile

                # Success or other status code
                return response

            except requests.RequestException as e:
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
        url = "https://www.linkedin.com/jobs/search/?currentJobId=4347732846&distance=25.0&geoId=105646813&keywords=%22scraping%22&origin=HISTORY"
        return self.make_request(url)

    def parse(self, html_content: str) -> list:
        soup = BeautifulSoup(html_content, "html.parser")
        job_list = soup.select("ul.jobs-search__results-list > li")
        self.logger.info(f"Found {len(job_list)} job postings on the page.")

        # if len(job_list) > 7:
        #     self.logger.info(
        #         f"Found {len(job_list)} job postings, saving HTML for review."
        #     )
        #     save_html(
        #         html_content,
        #         "linkedin_jobsearch_more_offers.html",
        #     )
        #     print()

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
                title = (
                    offer_soup.find("h1", class_="topcard__title").get_text().strip()
                )
                company = (
                    offer_soup.find("a", class_="topcard__org-name-link")
                    .get_text()
                    .strip()
                )
                location = (
                    offer_soup.find("span", class_="topcard__flavor--bullet")
                    .get_text()
                    .strip()
                )
                raw_date_posted_str = offer_soup.find("time")["datetime"].strip()
                date_posted_str = f"{raw_date_posted_str}T00:00:00Z"
                description = (
                    offer_soup.find("div", class_="description__text")
                    .get_text()
                    .strip()
                )
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
        response = self.linkedin_jobsearch_request()
        html_content = response.text
        jobs = self.parse(html_content)
        self.logger.info("Finished Linkedin scraping.")
        return jobs

    def scrape_test(self):
        html_content = load_html("./seed/linkedin_jobsearch_more_offers.html")
        return self.parse(html_content)
