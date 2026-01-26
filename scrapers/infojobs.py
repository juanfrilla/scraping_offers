import json
import random
import re

from bs4 import BeautifulSoup
from curl_cffi import requests

from constants import SEARCH_KEYWORDS
from logger import get_logger
from utils import determine_modality, find_keyword, normalize_string, read_json


class InfoJobsScraper:
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
        self.scraper_name = "Infojobs"
        self.logger = get_logger(self.scraper_name)

    def infojobs_jobsearch_request(self, keyword: str):
        url = "https://www.infojobs.net/jobsearch/search-results/list.xhtml"
        data = {
            "palabra": keyword,
            "normalizedJobTitleId": "",
            "of_provincia": "0",
            "canal": "0",
            "origen_busqueda": "0",
            "origen_accion": "0",
            "vieneUrlExecutive": "false",
        }
        return self.session.post(url, data=data, impersonate=self.get_impersonator())

    def get_impersonator(self):
        impersonate = random.choice(self.IMPERSONATE_LIST)
        self.logger.info(
            f"Selected Infojobs impersonate for scraping is {impersonate}."
        )
        return impersonate

    def convert_soup_to_json(self, soup: BeautifulSoup):
        target_script = None
        for script in soup.find_all("script"):
            if script.string and "window.__INITIAL_PROPS__" in script.string:
                target_script = script.string
                break

        if not target_script:
            self.logger.info("No script with __INITIAL_PROPS__ found")
            return {}
        match = re.search(r'JSON\.parse\("(.+)"\)', target_script, re.DOTALL)
        if match:
            raw = match.group(1)
            decoded = raw.encode("utf-8").decode("unicode_escape")
            data = json.loads(decoded)
            if "pde" in data:
                data["pde"] = json.loads(data["pde"])
            return data
        self.logger.info("No JSON.parse(...) match found")
        return {}

    def parse(self, json_data: dict) -> list:
        offers = json_data.get("offers", [])
        records = []
        for offer_idx, offer in enumerate(offers):
            self.logger.info(f"Parsing job posting {offer_idx + 1}/{len(offers)}")
            title = offer.get("title", "N/A")
            company = offer.get("companyName", "N/A")
            location = offer.get("city", "N/A")
            url = offer.get("link", "")
            date_posted_str = offer.get("publishedAt", "")
            description = offer.get("description", "")
            modality = offer.get("teleworking") or determine_modality(
                title, description
            )
            keyword_appeared = find_keyword(title, description)
            if keyword_appeared:
                records.append(
                    {
                        "title": title,
                        "company": normalize_string(company),
                        "location": self.determine_location(location),
                        "url": url,
                        "date_posted": date_posted_str,
                        "modality": self.determine_modality(modality),
                        "platform": "INFOJOBS",
                        "keyword_appeared": keyword_appeared,
                        "logo_url": offer.get("companyLogo", ""),
                    }
                )
        return records

    def determine_modality(self, modality: str) -> str:
        modality_map = {
            "Remoto": "Remote",
            "Híbrido": "Hybrid",
            "Presencial": "On-site",
            "Solo teletrabajo": "Remote",
        }
        return modality_map.get(modality, "N/A")

    def determine_location(self, location: str) -> str:
        locations = {"San Sebastián De Los Reyes": "Madrid", "Sabadell": "Barcelona"}
        location_lower = location.lower()
        normalized = normalize_string(location_lower)
        return locations.get(normalized, normalized)

    def scrape(self):
        all_jobs = []
        for keyword in SEARCH_KEYWORDS:
            response = self.infojobs_jobsearch_request(keyword)
            html = response.text
            soup = BeautifulSoup(html, "html.parser")
            json_data = self.convert_soup_to_json(soup)
            jobs = self.parse(json_data)
            all_jobs += jobs
        return all_jobs

    def scrape_test(self):
        json_data = read_json("./seed/infojobs_jobsearch.json")
        return self.parse(json_data)
