import json
import os
from datetime import date, datetime

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


def convert_datetime_str_to_other_format(
    entry_datetime: str, input_format: str, output_format: str
) -> str:
    dt = datetime.strptime(entry_datetime, input_format)
    return dt.strftime(output_format)
