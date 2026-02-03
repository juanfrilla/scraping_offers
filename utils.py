import html
import json
import os
import re
from datetime import date, datetime, timezone

from bs4 import BeautifulSoup

from constants import FORBIDDEN_KEYWORDS, LAST_SCRAPED_FILE

JSON_TYPES = ["application/json", "application/ld+json"]


def save_json(filename: str, data: dict):
    folder = os.path.dirname(filename)
    if folder:
        os.makedirs(folder, exist_ok=True)
    with open(filename, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=4, ensure_ascii=False)


def read_json(filename):
    with open(filename, "r") as f:
        info = json.loads(f.read())
        return info


def save_html(html_content, filename):
    with open(filename, "w", encoding="utf-8") as file:
        file.write(html_content)
    print(f"HTML saved to {filename}")


def load_html(filename):
    with open(filename, "r", encoding="utf-8") as file:
        html_content = file.read()
    return html_content


def last_scraped_today() -> bool:
    if not os.path.exists(LAST_SCRAPED_FILE):
        return False

    try:
        with open(LAST_SCRAPED_FILE, "r") as f:
            data = json.load(f)
    except (json.JSONDecodeError, IOError):
        return False

    last_date = data.get("date")
    return last_date == str(date.today())


def update_last_scraped() -> None:
    data = {"date": str(date.today())}

    with open(LAST_SCRAPED_FILE, "w") as f:
        json.dump(data, f, indent=2)


def normalize_string(s: str) -> str:
    return s.title().strip()


def determine_modality(title: str, description: str) -> str:
    text_to_search = f"{title} {description}".lower()
    keywords = {
        "Hybrid": ["hybrid", "híbrido", "hibrido", "mixto"],
        "Remote": ["remote", "remoto", "teletrabajo", "home office"],
        "On-site": ["on-site", "onsite", "presencial", "en oficina"],
    }
    for modality, terms in keywords.items():
        if any(term in text_to_search for term in terms):
            return modality

    return "N/A"


def find_keyword(title: str, text: str) -> str:
    if any(word in title.lower() for word in FORBIDDEN_KEYWORDS):
        return None
    _PATTERN = re.compile(r"\b(crawl|scrap|acqui)\w+", re.IGNORECASE)
    match = _PATTERN.search(text)

    return match.group(0) if match else None


def _json_from_ld(raw: str) -> dict:
    clean = raw.strip()
    clean = html.unescape(clean)
    return json.loads(clean)


def get_json_from_html(soup: BeautifulSoup) -> dict:
    for type in JSON_TYPES:
        script_tag = soup.find("script", type=type)
        if script_tag:
            raw_json_content = script_tag.string
            return _json_from_ld(raw_json_content)
    return {}


def parse_json_jobs_data(soup: BeautifulSoup) -> list:
    for type in JSON_TYPES:
        data_list = []
        scripts = soup.find_all("script", type=type) or []
        for script in scripts:
            json_data = _json_from_ld(script.string)
            data_list.append(json_data)
    return data_list


def from_timestamp_to_isoformat(ts_ms: int):
    dt = datetime.fromtimestamp(ts_ms / 1000, tz=timezone.utc)

    return dt.isoformat()


def filter_jobs(jobs: list) -> list:
    urls = set()
    filtered_jobs = []
    for job in jobs:
        url = job.get("url")
        if url not in urls:
            urls.add(url)
            filtered_jobs.append(job)
        else:
            continue
    return filtered_jobs
