from bs4 import BeautifulSoup
from curl_cffi import requests

from logger import get_logger
from utils import (
    find_keyword,
    get_json_from_html,
    normalize_string,
)


class RemoteRocketshipScraper:
    def __init__(self):
        self.session = requests.Session()
        self.scraper_name = "RemoteRocketship"
        self.logger = get_logger(self.scraper_name)

    def get_jobs_request(self):
        burp0_url = "https://www.remoterocketship.com/?page=1&sort=DateAdded&jobTitle=scraping%2Ccrawling%2Cscraper%2Ccrawler%2Cdata+acquisition%2Cdata+extraction&locations=Europe%2CSpain"

        response = self.session.get(burp0_url, impersonate="chrome131")
        soup = BeautifulSoup(response.text, "html.parser")
        return get_json_from_html(soup)

    def job_information_request(self, url: str):
        headers = {
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
            "accept-language": "es-ES,es;q=0.9",
            "cache-control": "no-cache",
            "priority": "u=0, i",
            "referer": "https://www.remoterocketship.com/?page=1&sort=DateAdded&jobTitle=scraping%2Ccrawling%2Cscraper%2Ccrawler%2Cdata+acquisition%2Cdata+extraction&locations=Europe%2CSpain",
            "sec-ch-ua": '"Google Chrome";v="131", "Chromium";v="131", "Not A(Brand";v="24"',
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": '"Windows"',
            "sec-fetch-dest": "document",
            "sec-fetch-mode": "navigate",
            "sec-fetch-site": "same-origin",
            "sec-fetch-user": "?1",
            "upgrade-insecure-requests": "1",
            "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
        }

        response = self.session.get(
            url,
            headers=headers,
            impersonate="chrome131",
        )
        soup = BeautifulSoup(response.text, "html.parser")
        return get_json_from_html(soup)

    def parse(self, response_json: dict) -> list:
        self.logger.info(f"response_json {response_json}")
        records = []
        jobs = (
            response_json.get("props", {})
            .get("pageProps", {})
            .get("initialJobOpenings", [])
        )

        for job_id, job in enumerate(jobs):
            self.logger.info(f"Parsing job posting {job_id + 1}/{len(jobs)}")
            company_info = job.get("company", {})
            company_name = company_info.get("name").replace(" ", "-")
            company_image_url = company_info.get("profilePicURL", "")
            slug = job.get("slug")
            url = (
                f"https://www.remoterocketship.com/company/{company_name}/jobs/{slug}/"
            )
            job_info = self.job_information_request(url)
            self.logger.info(f"job_info {job_info}")
            job_info_props = job_info.get("props", {}).get("pageProps", {})
            title = job.get("roleTitle")
            job_oppening = job_info_props.get("jobOpening")
            description = "\n".join(
                [
                    job_oppening.get("roleDescription"),
                    job_oppening.get("roleRequirements"),
                    job_oppening.get("benefits"),
                ]
            )

            location = job.get("location")
            date_posted_str = job.get("created_at")
            modality = "REMOTE"
            keyword_appeared = find_keyword(description)
            if keyword_appeared:
                records.append(
                    {
                        "title": title,
                        "company": normalize_string(company_name),
                        "location": location,
                        "url": job_oppening.get("url"),
                        "date_posted": date_posted_str,
                        "modality": modality,
                        "platform": "SIMPLYHIRED",
                        "keyword_appeared": keyword_appeared,
                        "logo_url": company_image_url,
                    }
                )

        return records

    def scrape(self):
        response_json = self.get_jobs_request()
        return self.parse(response_json)

        # pepe = self.job_information_request()

        # soup = BeautifulSoup(pepe.text, "html.parser")
        # json = get_json_from_html(soup)
        # print()
