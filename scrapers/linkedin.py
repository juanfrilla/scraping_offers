import html
import json

import curl_cffi as requests
from bs4 import BeautifulSoup

from logger import get_logger
from utils import determine_modality, keyword_counter, load_html, normalize_string


class LinkedinScraper:
    def __init__(self):
        self.logger = get_logger("linkedin_scraper")
        self.session = requests.Session()

    def linkedin_jobsearch_request(self):
        url = "https://www.linkedin.com/jobs/search/?currentJobId=4341337179&distance=25&geoId=105646813&keywords=%22scraping%22&origin=JOB_SEARCH_PAGE_QUERY_EXPANSION"
        burp0_headers = {
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
        }

        return self.session.get(url, headers=burp0_headers, impersonate="firefox135")

    def linkedin_entering_offer_request(self, offer_url: str):
        burp0_headers = {
            "Sec-Ch-Ua": '"Chromium";v="131", "Not A(Brand";v="24"',
            "Sec-Ch-Ua-Mobile": "?0",
            "Sec-Ch-Ua-Platform": '"Windows"',
            "Accept-Language": "es-ES,es;q=0.9",
            "Upgrade-Insecure-Requests": "1",
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
            "Sec-Fetch-Site": "none",
            "Sec-Fetch-Mode": "navigate",
            "Sec-Fetch-User": "?1",
            "Sec-Fetch-Dest": "document",
            "Accept-Encoding": "gzip, deflate, br",
            "Priority": "u=0, i",
        }
        return self.session.get(
            offer_url, headers=burp0_headers, impersonate="chrome131"
        )

    def parse(self, html_content: str) -> list:
        soup = BeautifulSoup(html_content, "html.parser")
        job_list = soup.select("ul.jobs-search__results-list > li")

        records = []

        for job_id, job in enumerate(job_list):
            url = job.select_one("a")["href"].strip()
            self.logger.info(f"Parsing job posting {job_id + 1}/{len(job_list)}")
            offer_response = self.linkedin_entering_offer_request(url)
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
        html_content = load_html("./seed/linkedin_jobsearch.html")
        return self.parse(html_content)
