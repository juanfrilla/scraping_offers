import os

import streamlit as st

from constants import DATA_FILE, LOG_DIR
from logger import get_logger
from scrapers.empleate import EmpleateScraper
from scrapers.infojobs import InfoJobsScraper
from scrapers.linkedin import LinkedinScraper
from utils import last_scraped_today, read_json, save_json, update_last_scraped


def scrape_everything():
    scrapers = [EmpleateScraper, LinkedinScraper, InfoJobsScraper]
    job_posts = []
    for ScraperClass in scrapers:
        scraper = ScraperClass()
        jobs = scraper.scrape()
        job_posts += jobs

    return job_posts


if __name__ == "__main__":
    os.makedirs(LOG_DIR, exist_ok=True)

    st.title("Daily Scraper")
    logger = get_logger("main")

    if last_scraped_today() and os.path.exists(DATA_FILE):
        logger.info("Data already scraped today. Loading from JSON...")
        jobs = read_json(DATA_FILE)
    else:
        logger.info("Scraping new data...")
        jobs = scrape_everything()
        save_json(DATA_FILE, jobs)
        update_last_scraped()
        logger.info("Data scraped and saved.")

    # Display jobs
    for job in jobs:
        with st.container():
            st.subheader(job.get("title", "No title"))
            st.write(f"**Company:** {job.get('company', 'N/A')}")
            st.write(f"**Location:** {job.get('location', 'N/A')}")
            st.write(f"**Posted:** {job.get('date_posted', 'N/A')}")
            st.write(f"**Platform:** {job.get('platform', 'N/A')}")
            st.write(f"[Apply Here]({job.get('url', '#')})")

        st.write("---")
