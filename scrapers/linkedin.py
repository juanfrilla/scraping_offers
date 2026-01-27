import random
import time
from datetime import datetime

import curl_cffi as requests
from bs4 import BeautifulSoup

from constants import FORBIDDEN_COMPANIES, FORBIDDEN_KEYWORDS, SEARCH_KEYWORDS
from logger import get_logger
from utils import (
    determine_modality,
    find_keyword,
    get_json_from_html,
    load_html,
    normalize_string,
)


class LinkedinScraper:
    def __init__(self):
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
        self.scraper_name = "Linkedin"
        self.logger = get_logger(self.scraper_name)

    def get_total_results(self, keyword: str, location_name: str, geo_id: str) -> int:
        """Obtiene el número total aproximado de resultados de la página principal."""
        search_url = f"https://www.linkedin.com/jobs/search?keywords={keyword}&location={location_name}&geoId={geo_id}&f_WT=2"
        response = self.make_request(search_url)
        if response and response.status_code == 200:
            soup = BeautifulSoup(response.text, "html.parser")
            count_tag = soup.select_one(".results-context-header__job-count")
            if count_tag:
                count_str = (
                    count_tag.get_text(strip=True)
                    .replace(",", "")
                    .replace("+", "")
                    .replace(".", "")
                )
                try:
                    return int(count_str)
                except ValueError:
                    return 0
        return 0

    def retrieve_offers(self) -> list:
        all_jobs = []
        self.impersonate = self.get_impersonator()

        locations = [
            {"name": "Spain", "url_name": "Espa%C3%B1a", "geoId": "105646813"},
            {"name": "Europe", "url_name": "Europa", "geoId": "100506914"},
            {
                "name": "United States",
                "url_name": "United%20States",
                "geoId": "103644278",
            },
        ]

        for location in locations:
            self.logger.info(f"--- Starting location: {location['name']} ---")
            for keyword in SEARCH_KEYWORDS:
                total_results = self.get_total_results(
                    keyword, location["url_name"], location["geoId"]
                )
                self.logger.info(
                    f"Total potential results for {keyword}: {total_results}"
                )
                max_fetch = min(total_results, 200) if total_results > 0 else 25
                start = 0

                while start < max_fetch:
                    self.logger.info(
                        f"Searching {keyword} in {location['name']} (Offset: {start})"
                    )

                    # 2. Construir URL de la API con el offset 'start'
                    self.job_search_url = (
                        f"https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?"
                        f"keywords=%22{keyword}%22&location={location['url_name']}&geoId={location['geoId']}"
                        f"&start={start}&f_WT=2"
                    )

                    response = self.linkedin_jobsearch_request()

                    if not response or response.status_code != 200:
                        break

                    soup = BeautifulSoup(response.text, "html.parser")
                    job_list = soup.select("li")

                    if not job_list:
                        self.logger.info("No more jobs found in this pagination.")
                        break

                    for job in job_list:
                        link_tag = job.select_one("a")
                        company_tag = job.select_one(".base-search-card__subtitle, h4")
                        if link_tag:
                            title = link_tag.get_text(strip=True).lower()
                            href = link_tag["href"].strip()
                            company_name = (
                                company_tag.get_text(strip=True).lower()
                                if company_tag
                                else ""
                            )

                            if any(
                                forbidden in title for forbidden in FORBIDDEN_KEYWORDS
                            ):
                                continue

                            if any(
                                company in company_name
                                for company in FORBIDDEN_COMPANIES
                            ):
                                self.logger.info(
                                    f"Skipping forbidden company: {company_name}"
                                )
                                continue

                            all_jobs.append(href)
                    start += 25
                    time.sleep(random.uniform(1.5, 3.0))

        return list(set(all_jobs))

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
        url = self.job_search_url
        self.logger.info(f"Scraping {url}")
        return self.make_request(url)

    def parse(self, urls: list) -> list:
        records = []
        external_ids = set()
        today_str = datetime.now().strftime("%Y-%m-%d")
        for job_id, url in enumerate(list(urls)):
            self.logger.info(f"Parsing job posting {job_id + 1}/{len(urls)}")
            self.logger.info(f"Parsing job {url}")
            offer_response = self.make_request(url)
            offer_soup = BeautifulSoup(offer_response.text, "html.parser")
            json_content = get_json_from_html(offer_soup)
            logo_url = ""
            if json_content:
                title = json_content.get("title", "N/A")
                company_data = json_content.get("hiringOrganization", {})
                company = company_data.get("name", "N/A")
                location = (
                    json_content.get("jobLocation", {})
                    .get("address", {})
                    .get("addressLocality", "N/A")
                )
                date_posted_str = json_content.get("datePosted", "")
                description = json_content.get("description", "")
                modality = determine_modality(title, description)
                logo_url = company_data.get("logo")
            else:
                raw_title = offer_soup.find("h1", class_="topcard__title")
                title = self.get_text(raw_title)
                raw_company = offer_soup.find("a", class_="topcard__org-name-link")
                company = self.get_text(raw_company)
                raw_location = offer_soup.find("span", class_="topcard__flavor--bullet")
                location = self.get_text(raw_location)
                try:
                    raw_date_posted_str = offer_soup.find("time")["datetime"].strip()
                except Exception:
                    self.logger.info(f"Erroneous date or not found for {url}")
                    raw_date_posted_str = today_str
                date_posted_str = f"{raw_date_posted_str}T00:00:00Z"
                raw_description = offer_soup.find("div", class_="description__text")
                description = self.get_text(raw_description)
                modality = determine_modality(title, description)
            determined_location = self.determine_location(location)
            external_id = f"{title}_{company}_{determined_location}"
            keyword_appeared = find_keyword(title, description)
            if keyword_appeared and external_id not in external_ids:
                external_ids.add(external_id)
                records.append(
                    {
                        "title": title,
                        "company": normalize_string(company),
                        "location": determined_location,
                        "url": url,
                        "date_posted": date_posted_str,
                        "modality": modality,
                        "platform": "LINKEDIN",
                        "keyword_appeared": keyword_appeared,
                        "logo_url": logo_url,
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
            "Greater Barcelona Metropolitan Area": "Barcelona",
            "Greater Madrid Metropolitan Area": "Madrid",
            "Vega Del Codorno": "Cuenca",
        }
        location_lower = location.lower()
        if location_lower.endswith("y alrededores"):
            location_lower = location_lower.replace("y alrededores", "")
        normalized = normalize_string(location_lower)
        return locations.get(normalized, normalized)

    def scrape(self):
        urls = self.retrieve_offers()
        jobs = self.parse(urls)
        return jobs

    def scrape_test(self):
        html_content = load_html("./seed/linkedin_jobsearch_more_offers.html")
        return self.parse(html_content)
