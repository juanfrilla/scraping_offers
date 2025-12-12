import json
import re

import curl_cffi as requests

from logger import get_logger
from utils import determine_modality, keyword_counter, normalize_string, read_json


class EmpleateScraper:
    def __init__(self):
        self.logger = get_logger("empleate_scraper")

    def empleate_jobsearch_request(self):
        url = "https://www.empleate.gob.es/empleate/open/offersearch/selectBuscador?q.op=AND&rows=10&sort=score%20desc&defType=edismax&df=titulo&facet=true&facet.field=paisF&facet.field=provinciaF&facet.field=provincia&facet.field=categoria&facet.field=categoriaF&facet.field=subcategoriaF&facet.field=subcategoria&facet.field=origen&facet.field=tipoContrato&facet.field=tipoContratoN&facet.field=noMeInteresa&facet.field=educacionF&facet.field=fechaCreacionPortal&facet.field=jornadaF&facet.field=experienciaF&facet.field=educacion&facet.field=minExperiencia&facet.field=jornada&facet.field=pais&facet.field=discapacidad&facet.field=cno&facet.field=portales&facet.field=showPortalPu&facet.field=showPortalPr&facet.mincount=1&f.topics.facet.limit=50&json.nl=map&fq=(speStateId%3A1%20OR%20speStateId%3A4)%20AND%20checkVisible%3A1&fl=*%2C%20score&q=%22scraping%22&wt=json&json.wrf=jQuery110204302387032510525_1765140545478&_=1765140545479"
        headers = {
            "Sec-Ch-Ua-Platform": '"Windows"',
            "X-Requested-With": "XMLHttpRequest",
            "Accept-Language": "es-ES,es;q=0.9",
            "Accept": "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01",
            "Sec-Ch-Ua": '"Not_A Brand";v="99", "Chromium";v="131"',
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
            "Sec-Ch-Ua-Mobile": "?0",
            "Sec-Fetch-Site": "same-origin",
            "Sec-Fetch-Mode": "cors",
            "Sec-Fetch-Dest": "empty",
            "Referer": "https://www.empleate.gob.es/empleo/",
            "Accept-Encoding": "gzip, deflate, br",
            "Priority": "u=1, i",
            "Connection": "keep-alive",
        }
        return requests.get(url, headers=headers, impersonate="chrome131", verify=False)

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
        self.logger.info(f"Found {len(docs)} job postings.")
        for doc_id, doc in enumerate(docs):
            self.logger.info(f"Parsing job posting {doc_id + 1}/{len(docs)}")
            title = doc.get("titulo", "N/A")
            company = doc.get("contacto", "N/A")
            location = doc.get("provinciaF", "N/A")
            url = doc.get("url", "")
            date_posted_str = doc.get("fechaCreacionPortal", "")
            platform = doc.get("entitytype")
            description = doc.get("contenido", "")
            records.append(
                {
                    "title": title,
                    "company": normalize_string(company),
                    "location": normalize_string(location),
                    "url": url,
                    "date_posted": date_posted_str,
                    "platform": platform,
                    "modality": determine_modality(title, doc.get("contenido", "")),
                    "keywords": keyword_counter(description),
                }
            )
        return records

    def scrape(self):
        self.logger.info("Starting Empleate scraping.")
        response = self.empleate_jobsearch_request()
        content = response.content
        json_data = self.convert_content_to_json(content)
        jobs = self.parse(json_data)
        self.logger.info("Finished Empleate scraping.")
        return jobs

    def scrape_test(self):
        json_data = read_json("./seed/empleate_jobsearch.json")
        return self.parse(json_data)
