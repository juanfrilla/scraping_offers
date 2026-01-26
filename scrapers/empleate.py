import json
import re

import curl_cffi as requests

from logger import get_logger
from utils import (
    determine_modality,
    filter_jobs,
    find_keyword,
    normalize_string,
    read_json,
)


class EmpleateScraper:
    def __init__(self):
        self.scraper_name = "Empleate"
        self.logger = get_logger(self.scraper_name)

    def empleate_jobsearch_request(self, keyword: str):
        url = f"https://www.empleate.gob.es/empleate/open/offersearch/selectBuscador?q.op=AND&rows=10&sort=score%20desc&defType=edismax&df=titulo&facet=true&facet.field=paisF&facet.field=provinciaF&facet.field=provincia&facet.field=categoria&facet.field=categoriaF&facet.field=subcategoriaF&facet.field=subcategoria&facet.field=origen&facet.field=tipoContrato&facet.field=tipoContratoN&facet.field=noMeInteresa&facet.field=educacionF&facet.field=fechaCreacionPortal&facet.field=jornadaF&facet.field=experienciaF&facet.field=educacion&facet.field=minExperiencia&facet.field=jornada&facet.field=pais&facet.field=discapacidad&facet.field=cno&facet.field=portales&facet.field=showPortalPu&facet.field=showPortalPr&facet.mincount=1&f.topics.facet.limit=50&json.nl=map&fq=(speStateId%3A1%20OR%20speStateId%3A4)%20AND%20checkVisible%3A1&fl=*%2C%20score&q={keyword}&wt=json&json.wrf=jQuery110200911569443944088_1766487062418&_=1766487062428"
        return requests.get(url, impersonate="chrome131", verify=False)

    def convert_content_to_json(self, content: bytes) -> dict:
        text = content.decode("utf-8")
        match = re.search(r"\((.*)\)$", text)
        if not match:
            self.logger.info("JSONP format not recognized.")
            return {}
        json_str = match.group(1)
        return json.loads(json_str)

    def parse(self, json_data: dict) -> list:
        records = []
        docs = json_data.get("response", {}).get("docs", [])
        for doc_id, doc in enumerate(docs):
            self.logger.info(f"Parsing job posting {doc_id + 1}/{len(docs)}")
            title = doc.get("titulo", "N/A")
            company = doc.get("contacto", "N/A")
            location = doc.get("provinciaF", "N/A")
            url = doc.get("url", "")
            date_posted_str = doc.get("fechaCreacionPortal", "")
            platform = doc.get("entitytype")
            description = doc.get("contenido", "")
            keyword_appeared = find_keyword(description, title)
            if keyword_appeared:
                records.append(
                    {
                        "title": title,
                        "company": normalize_string(company),
                        "location": normalize_string(location),
                        "url": url,
                        "date_posted": date_posted_str,
                        "platform": platform,
                        "modality": determine_modality(title, doc.get("contenido", "")),
                        "keyword_appeared": keyword_appeared,
                        "logo_url": "",
                    }
                )
        return records

    def scrape(self):
        keywords = [
            "scraping",
            "crawling",
            "data%20aquisition",
            "data%20extraction",
            "scraper",
            "crawler",
        ]
        all_jobs = []
        for keyword in keywords:
            response = self.empleate_jobsearch_request(keyword)
            content = response.content
            json_data = self.convert_content_to_json(content)
            jobs = self.parse(json_data)
            all_jobs += jobs
        return filter_jobs(all_jobs)

    def scrape_test(self):
        json_data = read_json("./seed/empleate_jobsearch.json")
        return self.parse(json_data)
