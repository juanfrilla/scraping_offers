import os
from datetime import date, datetime

import streamlit as st

from constants import DATA_FILE, LOG_DIR
from logger import get_logger
from scrapers.empleate import EmpleateScraper
from scrapers.infojobs import InfoJobsScraper
from scrapers.linkedin import LinkedinScraper
from utils import (
    convert_datetime_str_to_other_format,
    last_scraped_today,
    read_json,
    save_json,
    update_last_scraped,
)


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

    locations = sorted(list({job.get("location", "N/A") for job in jobs}))
    companies = sorted(list({job.get("company", "N/A") for job in jobs}))
    platforms = sorted(list({job.get("platform", "N/A") for job in jobs}))

    st.sidebar.header("Filters")
    selected_location = st.sidebar.selectbox("Location", ["All"] + locations)
    selected_company = st.sidebar.selectbox("Company", ["All"] + companies)
    selected_platform = st.sidebar.selectbox("Platform", ["All"] + platforms)
    job_dates = [
        datetime.strptime(job["date_posted"], "%Y-%m-%dT%H:%M:%SZ").date()
        for job in jobs
    ]
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
            or (
                start_date
                <= datetime.strptime(
                    job.get("date_posted"), "%Y-%m-%dT%H:%M:%SZ"
                ).date()
                <= end_date
            )
        )
    ]

    sorted_jobs_by_datetime = sorted(
        filtered_jobs,
        key=lambda x: datetime.strptime(x["date_posted"], "%Y-%m-%dT%H:%M:%SZ"),
        reverse=True,
    )
    st.title(f"Web Scraping Offers in Spain ({len(sorted_jobs_by_datetime)} results)")

    for job in sorted_jobs_by_datetime:
        with st.container():
            st.subheader(job.get("title", "No title"))

            col1, col2 = st.columns([3, 1])

            date_posted = convert_datetime_str_to_other_format(
                job.get("date_posted", "N/A"),
                "%Y-%m-%dT%H:%M:%SZ",
                "%d/%m/%Y",
            )

            with col1:
                st.write(f"🏢**Company:** {job.get('company', 'N/A')}")
                st.write(f"📍**Location:** {job.get('location', 'N/A')}")
                st.write(f"📅**Posted:** {date_posted}")
                st.write(f"💻**Platform:** {job.get('platform', 'N/A')}")

            with col2:
                st.markdown(
                    f"[Apply Here]({job.get('url', '#')})", unsafe_allow_html=True
                )

        st.write("---")
