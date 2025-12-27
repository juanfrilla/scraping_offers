import curl_cffi as requests
from bs4 import BeautifulSoup

from utils import get_json_from_html


def get_proxies_request():
    headers = {
        "accept": "application/json, text/plain, */*",
        "accept-language": "es-ES,es;q=0.9",
        "origin": "https://proxyscrape.com",
        "priority": "u=1, i",
        "referer": "https://proxyscrape.com/",
        "sec-ch-ua": '"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"',
        "sec-ch-ua-mobile": "?0",
        "sec-ch-ua-platform": '"Windows"',
        "sec-fetch-dest": "empty",
        "sec-fetch-mode": "cors",
        "sec-fetch-site": "same-site",
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
    }

    params = {
        "request": "get_proxies",
        "skip": "0",
        "proxy_format": "protocolipport",
        "format": "json",
        "limit": "15",
    }

    response = requests.get(
        "https://api.proxyscrape.com/v4/free-proxy-list/get",
        params=params,
        headers=headers,
    )
    return response.json()


if __name__ == "__main__":
    rjson = get_proxies_request()
    proxies = rjson.get("proxies")
    for proxy_data in proxies:
        burp0_url = "https://www.simplyhired.es/search?q=%22crawler%22+or+%22scraping%22&l=espa%25c3%25b1a"
        proxy = proxy_data.get("proxy")
        print(f"using proxy {proxy}")
        try:
            response = requests.get(
                burp0_url,
                impersonate="chrome131",
                proxies={"http": proxy, "https": proxy},
            )
        except Exception:
            print("exception")
            continue
        html_json = get_json_from_html(BeautifulSoup(response.text))
        if html_json != {}:
            print()
        else:
            continue
