from datetime import datetime

from bs4 import BeautifulSoup

from utils import load_html


def test_linkedin_jobsearch_html_saved():
    html_content = load_html("linkedin_jobsearch.html")
    soup = BeautifulSoup(html_content, "html.parser")
    job_list = soup.select("ul.jobs-search__results-list > li")
    for job in job_list:
        title = job.select("h3")[0].text.strip()
        company = job.select("h4")[0].text.strip()
        location = job.select("span.job-search-card__location")[0].text.strip()
        url = job.select("a")[0]["href"].strip()
        date_posted = datetime.strptime(
            job.select("time")[0]["datetime"].strip(), "%Y-%m-%d"
        ).date()
        print(f"Title: {title}")
        print(f"Company: {company}")
        print(f"Location: {location}")
        print(f"URL: {url}")
        print(f"Date Posted: {date_posted}")


if __name__ == "__main__":
    test_linkedin_jobsearch_html_saved()
