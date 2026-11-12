import requests

burp0_url = "https://www.remoterocketship.com:443/remote-jobs/?page=1&sort=DateAdded&jobTitle=scraping%2Ccrawling%2Cscraper%2Ccrawler%2Cdata%2520acquisition&locations=Europe%2CSpain"
burp0_headers = {
    "Sec-Ch-Ua": '"Not:A-Brand";v="8", "Chromium";v="156"',
    "Sec-Ch-Ua-Mobile": "?0",
    "Sec-Ch-Ua-Platform": '"Windows"',
    "Accept-Language": "es-ES,es;q=0.9",
    "Upgrade-Insecure-Requests": "1",
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/jxl,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
    "Sec-Fetch-Site": "none",
    "Sec-Fetch-Mode": "navigate",
    "Sec-Fetch-User": "?1",
    "Sec-Fetch-Dest": "document",
    "Accept-Encoding": "gzip, deflate, br",
    "Priority": "u=0, i",
    "Connection": "keep-alive",
}
requests.get(burp0_url, headers=burp0_headers)
