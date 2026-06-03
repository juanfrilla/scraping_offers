import os
import subprocess
import threading
from datetime import date, datetime, timezone

import streamlit as st
from dateutil.parser import isoparse
from dotenv import load_dotenv

from py.constants import DATA_FILE, LOG_DIR
from py.logger import get_logger
from py.utils import (
    last_scraped_today,
    read_json,
    update_last_scraped,
)

load_dotenv()


def run_go_scraper() -> tuple[bool, str]:
    ENV = os.getenv("ENV", "prod")
    logger = get_logger("main")
    env = os.environ.copy()
    command = (
        ["go", "run", "./go/main.go"] if ENV == "local" else ["./go/bin/scraper_go"]
    )
    logger.info(f"Executing {command}")

    try:
        process = subprocess.Popen(
            command,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1,
            env=env,
        )

        logs = []
        err_logs = []

        def read_stdout():
            for line in process.stdout:
                line = line.rstrip()
                logs.append(line)
                logger.info(f"[GO] {line}")

        def read_stderr():
            for line in process.stderr:
                line = line.rstrip()
                err_logs.append(line)
                logger.error(f"[GO-ERR] {line}")

        t1 = threading.Thread(target=read_stdout)
        t2 = threading.Thread(target=read_stderr)

        t1.start()
        t2.start()

        process.wait(timeout=120)

        t1.join()
        t2.join()

        if process.returncode != 0:
            return False, "\n".join(err_logs) or "\n".join(logs)

        return True, "\n".join(logs)

    except subprocess.TimeoutExpired:
        process.kill()
        return False, "Timeout: scraper tardó más de 120s."

    except FileNotFoundError:
        return False, "Binary './scraper_go' not found."


@st.cache_resource
def get_scraper_lock():
    return threading.Lock()


@st.cache_data(ttl=86400, show_spinner=False)
def get_jobs_data() -> list:
    """Obtiene los jobs: del caché JSON o ejecutando el scraper de Go."""
    if last_scraped_today() and os.path.exists(DATA_FILE):
        return read_json(DATA_FILE)
    success, message = run_go_scraper()

    if success:
        update_last_scraped()
        return read_json(DATA_FILE)
    else:
        st.error(f"Error en el scraper: {message}")
        if os.path.exists(DATA_FILE):
            st.warning("Usando datos del último scraping disponible.")
            return read_json(DATA_FILE)
        return []


with st.spinner("Scraping jobs..."):
    os.makedirs(LOG_DIR, exist_ok=True)
    jobs = get_jobs_data()

    st.sidebar.header("Filters")
    platforms = sorted(list({job.get("platform", "N/A") for job in jobs}))
    selected_platform = st.sidebar.selectbox("Platform", ["All"] + platforms)

    jobs_after_platform = [
        j
        for j in jobs
        if selected_platform == "All" or j.get("platform") == selected_platform
    ]
    companies = sorted(list({j.get("company", "N/A") for j in jobs_after_platform}))
    selected_company = st.sidebar.selectbox("Company", ["All"] + companies)

    jobs_after_company = [
        j
        for j in jobs_after_platform
        if selected_company == "All" or j.get("company") == selected_company
    ]
    locations = sorted(list({j.get("location", "N/A") for j in jobs_after_company}))
    selected_location = st.sidebar.selectbox("Location", ["All"] + locations)

    jobs_after_location = [
        j
        for j in jobs_after_company
        if selected_location == "All" or j.get("location") == selected_location
    ]
    modalities = sorted(list({j.get("modality", "N/A") for j in jobs_after_location}))
    selected_modality = st.sidebar.selectbox("Modality", ["All"] + modalities)

    job_dates = [isoparse(job["date_posted"]).date() for job in jobs]
    min_date = min(job_dates) if job_dates else date.today()
    max_date = max(job_dates) if job_dates else date.today()

    date_range = st.sidebar.date_input(
        "Filter by date posted (range)",
        min_value=min_date,
        max_value=max_date,
        value=(min_date, max_date),
        format="DD/MM/YYYY",
    )

    if isinstance(date_range, tuple) and len(date_range) == 2:
        start_date, end_date = date_range
    else:
        start_date = end_date = None

    filtered_jobs = [
        job
        for job in jobs_after_location
        if (selected_modality == "All" or job.get("modality") == selected_modality)
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
    st.title(
        f"Web Scraping Offers in or from Spain ({len(sorted_jobs_by_datetime)} results)"
    )

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
            st.subheader(f"{title} 🆕" if is_new else title)
            col_logo, col_info, col_button = st.columns([1, 4, 1.5])

            with col_logo:
                logo_url = job.get("logo_url")
                if logo_url:
                    st.image(logo_url, width=80)
                else:
                    st.write("🏢")

            with col_info:
                st.write(f"**{job.get('company', 'N/A')}**")
                st.write(
                    f"📍 {job.get('location', 'N/A')} | 💼 {job.get('modality', 'N/A')}"
                )
                st.write(f"📅 Posted: {date_posted} | 🌐 {job.get('platform', 'N/A')}")
                st.write("**Keyword appeared:** " + job.get("keyword_appeared", "N/A"))

            with col_button:
                st.markdown(
                    f"""<a href="{job.get("url", "#")}" target="_blank">
                        <button style="width:100%; border-radius:5px; background-color:#007BFF; color:white; border:none; padding:10px;">
                            Apply Here
                        </button>
                    </a>""",
                    unsafe_allow_html=True,
                )

        st.write("---")
