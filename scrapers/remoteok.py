import curl_cffi as requests
from bs4 import BeautifulSoup

from logger import get_logger
from utils import find_keyword, get_json_from_html, normalize_string


class RemoteOKScraper:
    def __init__(self):
        self.session = requests.Session()
        self.scraper_name = "RemoteOK"
        self.logger = get_logger(self.scraper_name)

    def get_jobs_request(self, keyword: str):
        burp0_url = f"https://remoteok.com:443/?tags={keyword}&action=get_jobs&premium=0&regular=1"
        return self.session.get(burp0_url, impersonate="chrome131")

    def parse(self, soup: BeautifulSoup) -> list:
        records = []
        jobs = soup.select("tr.job")
        # jobs = parse_json_jobs_data(soup)
        for job_id, job in enumerate(jobs):
            data_id = job.get("data-id")
            expanded_description = soup.select(f"tr.expand.expand-{data_id}")[0].text
            job_json = get_json_from_html(job)
            self.logger.info(f"Parsing job posting {job_id + 1}/{len(jobs)}")
            bot_url = job.get("data-url")
            title = job_json.get("title", "N/A")
            company = job_json.get("hiringOrganization", {}).get("name")
            raw_date_posted = job_json.get("datePosted")
            modality = "REMOTE"
            location = ",".join(
                [
                    loc.get("name", "")
                    for loc in job_json.get("applicantLocationRequirements", [])
                ]
            )
            keyword_appeared = find_keyword(expanded_description)
            if keyword_appeared:
                records.append(
                    {
                        "title": title,
                        "company": normalize_string(company),
                        "location": location,
                        "url": f"https://www.remoteok.com{bot_url}",
                        "date_posted": raw_date_posted,
                        "modality": modality,
                        "platform": "REMOTEOK",
                        "keyword_appeared": keyword_appeared,
                        "logo_url": job_json.get("image"),
                    }
                )

        return records

    def scrape(self):
        keywords = [
            "python",
            "javascript",
        ]
        all_jobs = []
        for keyword in keywords:
            self.logger.info(f"Scraping word {keyword}")
            jobs_response = self.get_jobs_request(keyword)
            soup = BeautifulSoup(jobs_response.text, "html.parser")
            jobs = self.parse(soup)
            self.logger.info(f"Retrieved {len(jobs)} for {keyword}")
            all_jobs += jobs
        return all_jobs