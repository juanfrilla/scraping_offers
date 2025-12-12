import json
import re

from bs4 import BeautifulSoup
from curl_cffi import requests

from logger import get_logger
from utils import read_json


class InfoJobsScraper:
    def __init__(self):
        self.logger = get_logger("infojobs_scraper")

    def infojobs_jobsearch_request(self):
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
            "palabra": "scraping",
            "normalizedJobTitleId": "",
            "of_provincia": "0",
            "canal": "0",
            "origen_busqueda": "0",
            "origen_accion": "0",
            "vieneUrlExecutive": "false",
        }
        return requests.post(url, headers=headers, data=data, impersonate="chrome131")

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
            title = offer.get("title", "")
            company = offer.get("companyName", "")
            location = offer.get("city", "")
            url = offer.get("link", "")
            date_posted_str = offer.get("publishedAt", "")
            records.append(
                {
                    "title": title,
                    "company": company,
                    "location": location,
                    "url": url,
                    "date_posted": date_posted_str,
                    "platform": "INFOJOBS",
                }
            )
        return records

    def scrape(self):
        self.logger.info("Starting Infojobs scraping.")
        response = self.infojobs_jobsearch_request()
        html = response.text
        soup = BeautifulSoup(html, "html.parser")
        json_data = self.convert_soup_to_json(soup)
        jobs = self.parse(json_data)
        self.logger.info("Finished Infojobs scraping.")
        return jobs

    def scrape_test(self):
        json_data = read_json("infojobs_jobsearch.json")
        return self.parse(json_data)
