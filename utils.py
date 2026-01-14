import html
import json
import os
import re
from datetime import date, datetime, timezone

from bs4 import BeautifulSoup

from constants import LAST_SCRAPED_FILE


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
    title_lower = title.lower()
    description_lower = description.lower()

    if (
        "remote" in title_lower
        or "remote" in description_lower
        or "remoto" in title_lower
        or "remoto" in description_lower
    ):
        return "Remote"
    elif (
        "hybrid" in title_lower
        or "hybrid" in description_lower
        or "híbrido" in title_lower
        or "híbrido" in description_lower
    ):
        return "Hybrid"
    elif (
        "on-site" in title_lower
        or "on-site" in description_lower
        or "presencial" in title_lower
        or "presencial" in description_lower
    ):
        return "On-site"
    else:
        return "N/A"


def find_keyword(text: str) -> str:
    _PATTERN = re.compile(r"\b(crawl|scrap|acqui|extract)\w*", re.IGNORECASE)

    match = _PATTERN.search(text)
    return match.group(0) if match else None


def json_from_ld(raw: str) -> dict:
    clean = raw.strip()
    clean = html.unescape(clean)
    return json.loads(clean)


def get_json_from_html(soup: BeautifulSoup) -> dict:
    types = ["application/json", "application/ld+json"]
    for type in types:
        script_tag = soup.find("script", type=type)
        if script_tag:
            raw_json_content = script_tag.string
            return json_from_ld(raw_json_content)
    return {}


def from_timestamp_to_isoformat(ts_ms: int):
    dt = datetime.fromtimestamp(ts_ms / 1000, tz=timezone.utc)

    return dt.isoformat()
