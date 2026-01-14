import json
import re

from bs4 import BeautifulSoup
from curl_cffi import requests

from logger import get_logger
from utils import determine_modality, find_keyword, normalize_string, read_json


class InfoJobsScraper:
    def __init__(self):
        self.logger = get_logger("infojobs_scraper")
        self.session = requests.Session()

    def infojobs_jobsearch_request(self, keyword: str):
        url = "https://www.infojobs.net/jobsearch/search-results/list.xhtml"
        headers = {
            "Cache-Control": "max-age=0",
            "Sec-Ch-Ua": '"Not_A Brand";v="99", "Chromium";v="131"',
            "Sec-Ch-Ua-Mobile": "?0",
            "Sec-Ch-Ua-Platform": '"Windows"',
            "Accept-Language": "es-ES,es;q=0.9",
            "Origin": "https://www.infojobs.net",
            "Content-Type": "application/x-www-form-urlencoded",
            "Upgrade-Insecure-Requests": "1",
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
            "Sec-Fetch-Site": "same-origin",
            "Sec-Fetch-Mode": "navigate",
            "Sec-Fetch-User": "?1",
            "Sec-Fetch-Dest": "document",
            "Referer": "https://www.infojobs.net/?nocache=true",
            "Accept-Encoding": "gzip, deflate, br",
            "Priority": "u=0, i",
        }
        data = {
            "palabra": keyword,
            "normalizedJobTitleId": "",
            "of_provincia": "0",
            "canal": "0",
            "origen_busqueda": "0",
            "origen_accion": "0",
            "vieneUrlExecutive": "false",
        }
        return self.session.post(url, data=data, impersonate="chrome136")

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
            keyword_appeared = find_keyword(description)
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
        self.logger.info("Starting Infojobs scraping.")
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
            response = self.infojobs_jobsearch_request(keyword)
            html = response.text
            soup = BeautifulSoup(html, "html.parser")
            json_data = self.convert_soup_to_json(soup)
            jobs = self.parse(json_data)
            self.logger.info(f"Retrieved {len(jobs)} for {keyword}")
            all_jobs += jobs
            self.logger.info("Finished Infojobs scraping.")
        return all_jobs

    def scrape_test(self):
        json_data = read_json("./seed/infojobs_jobsearch.json")
        return self.parse(json_data)
