import json
import os
from datetime import date

from py.constants import LAST_SCRAPED_FILE


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
