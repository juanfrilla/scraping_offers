from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

LOG_DIR = BASE_DIR / "logs"
DATA_FILE = BASE_DIR / "data" / "jobs.json"
LAST_SCRAPED_FILE = BASE_DIR / "data" / "last_scraped.json"
