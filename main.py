import os
import threading
from datetime import date, datetime, timezone

import streamlit as st
from dateutil.parser import isoparse
from dotenv import load_dotenv

from constants import DATA_FILE, LOG_DIR
from logger import get_logger
from scrapers.empleate import EmpleateScraper
from scrapers.infojobs import InfoJobsScraper
from scrapers.linkedin import LinkedinScraper
from scrapers.simplyhired import SimplyHiredScraper
from utils import (
    last_scraped_today,
    read_json,
    save_json,
    update_last_scraped,
)

load_dotenv()


def scrape_everything():
    ENV = os.getenv("ENV", "server")
    in_server = ENV == "server"
    if in_server:
        scrapers = [LinkedinScraper, InfoJobsScraper, EmpleateScraper]
    else:
        scrapers = [
            SimplyHiredScraper,
            LinkedinScraper,
            InfoJobsScraper,
            EmpleateScraper,
        ]
    job_posts = []
    for ScraperClass in scrapers:
        scraper = ScraperClass()
        jobs = scraper.scrape()
        job_posts += jobs

    return job_posts


@st.cache_resource
def get_scraper_lock():
    return threading.Lock()


@st.cache_data(ttl=86400, show_spinner=False)
def get_jobs_data():
    lock = get_scraper_lock()
    with lock:
        if last_scraped_today() and os.path.exists(DATA_FILE):
            return read_json(DATA_FILE)
        jobs = scrape_everything()
        save_json(DATA_FILE, jobs)
        update_last_scraped()
        return jobs


with st.spinner("Scraping jobs..."):
    os.makedirs(LOG_DIR, exist_ok=True)
    logger = get_logger("main")
    jobs = get_jobs_data()

    locations = sorted(list({job.get("location", "N/A") for job in jobs}))
    companies = sorted(list({job.get("company", "N/A") for job in jobs}))
    platforms = sorted(list({job.get("platform", "N/A") for job in jobs}))
    modalities = sorted(list({job.get("modality", "N/A") for job in jobs}))

    st.sidebar.header("Filters")
    selected_location = st.sidebar.selectbox("Location", ["All"] + locations)
    selected_company = st.sidebar.selectbox("Company", ["All"] + companies)
    selected_platform = st.sidebar.selectbox("Platform", ["All"] + platforms)
    selected_modality = st.sidebar.selectbox("Modality", ["All"] + modalities)
    job_dates = [isoparse(job["date_posted"]).date() for job in jobs]
    min_date = min(job_dates) if job_dates else date.today()
    max_date = max(job_dates) if job_dates else date.today()
    date_range = st.sidebar.date_input(
        "Filter by date posted (range)",
        min_value=min_date,
        max_value=max_date,
        value=(min_date, max_date),
        help="Select a date range to show jobs posted within it",
        format="DD/MM/YYYY",
    )
    if isinstance(date_range, tuple) and len(date_range) == 2:
        start_date, end_date = date_range
    else:
        start_date = end_date = None

    filtered_jobs = [
        job
        for job in jobs
        if (selected_location == "All" or job.get("location") == selected_location)
        and (selected_company == "All" or job.get("company") == selected_company)
        and (selected_platform == "All" or job.get("platform") == selected_platform)
        and (
            (not start_date or not end_date)
            or (start_date <= isoparse(job["date_posted"]).date() <= end_date)
        )
    ]

    sorted_jobs_by_datetime = sorted(
        filtered_jobs,
        key=lambda x: isoparse(x["date_posted"]),
        reverse=True,
    )
    st.title(f"Web Scraping Offers in Spain ({len(sorted_jobs_by_datetime)} results)")

    for job in sorted_jobs_by_datetime:
        raw_date = job.get("date_posted")

        is_new = False
        date_posted = "N/A"

        if raw_date:
            posted_dt = isoparse(raw_date)
            now = datetime.now(timezone.utc)

        is_new = (now - posted_dt).days < 3
        date_posted = posted_dt.strftime("%d/%m/%Y")
        with st.container():
            title = job.get("title", "No title")
            date_posted = posted_dt.strftime("%d/%m/%Y")

            if is_new:
                st.subheader(f"{title} 🆕")
            else:
                st.subheader(title)

            col1, col2 = st.columns([3, 1])

            with col1:
                st.write(f"🏢**Company:** {job.get('company', 'N/A')}")
                st.write(f"📍**Location:** {job.get('location', 'N/A')}")
                st.write(f"💼**Modality:** {job.get('modality', 'N/A')}")
                st.write(f"📅**Posted:** {date_posted}")
                st.write(f"🌐**Platform:** {job.get('platform', 'N/A')}")
                st.write("**Keyword appeared:** " + job.get("keyword_appeared", []))

            with col2:
                st.markdown(
                    f"[Apply Here]({job.get('url', '#')})", unsafe_allow_html=True
                )

        st.write("---")
