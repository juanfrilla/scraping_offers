import json
import os
import re
from collections import Counter
from datetime import date

from dateutil.parser import isoparse

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


def convert_isodatetime_str_to_other_format(
    entry_datetime: str, output_format: str
) -> str:
    dt = isoparse(entry_datetime)
    return dt.strftime(output_format)


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


def keyword_counter(text, min_length=5):
    text = text.lower()
    words = re.findall(r"[a-zA-Z]+", text)
    filtered = [w for w in words if len(w) >= min_length]
    counts = Counter(filtered)
    sorted_keywords = sorted(counts.items(), key=lambda x: x[1], reverse=True)
    return [keyword for keyword, _ in sorted_keywords[:4]]
