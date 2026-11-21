from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

LOG_DIR = BASE_DIR / "logs"
LAST_SCRAPED_FILE = BASE_DIR / "data" / "last_scraped.json"
DATABASE_FILE = BASE_DIR / "data" / "jobs.db"
