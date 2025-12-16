from bs4 import BeautifulSoup
from curl_cffi import requests

profiles = [
    # =========================
    # FIREFOX
    # =========================
    {
        "impersonate": "firefox133",
        "headers": {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.9",
            "Accept-Encoding": "gzip, deflate, br, zstd",
            "Upgrade-Insecure-Requests": "1",
            "Sec-Fetch-Site": "none",
            "Sec-Fetch-Mode": "navigate",
            "Sec-Fetch-User": "?1",
            "Sec-Fetch-Dest": "document",
        },
    },
    # =========================
    # TOR
    # =========================
    {
        "impersonate": "tor145",
        "headers": {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; rv:115.0) Gecko/20100101 Firefox/115.0",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.5",
            "Accept-Encoding": "gzip, deflate, br",
            "Upgrade-Insecure-Requests": "1",
        },
    },
    # =========================
    # CHROME (DESKTOP)
    # =========================
    *[
        {
            "impersonate": name,
            "headers": {
                "User-Agent": f"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/{ver}.0.0.0 Safari/537.36",
                "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
                "Accept-Language": "en-US,en;q=0.9",
                "Accept-Encoding": "gzip, deflate, br, zstd",
                "Upgrade-Insecure-Requests": "1",
                "Sec-Fetch-Site": "none",
                "Sec-Fetch-Mode": "navigate",
                "Sec-Fetch-User": "?1",
                "Sec-Fetch-Dest": "document",
                "Sec-CH-UA": '"Chromium";v="{}","Google Chrome";v="{}","Not=A?Brand";v="99"'.format(
                    ver, ver
                ),
                "Sec-CH-UA-Mobile": "?0",
                "Sec-CH-UA-Platform": '"Windows"',
            },
        }
        for name, ver in [
            ("chrome99", 99),
            ("chrome100", 100),
            ("chrome101", 101),
            ("chrome104", 104),
            ("chrome107", 107),
            ("chrome110", 110),
            ("chrome116", 116),
            ("chrome119", 119),
            ("chrome120", 120),
            ("chrome123", 123),
            ("chrome124", 124),
            ("chrome131", 131),
            ("chrome133a", 133),
            ("chrome136", 136),
        ]
    ],
    # =========================
    # CHROME (ANDROID)
    # =========================
    {
        "impersonate": "chrome99_android",
        "headers": {
            "User-Agent": "Mozilla/5.0 (Linux; Android 12; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.73 Mobile Safari/537.36",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.9",
            "Accept-Encoding": "gzip, deflate, br",
        },
    },
    {
        "impersonate": "chrome131_android",
        "headers": {
            "User-Agent": "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.9",
            "Accept-Encoding": "gzip, deflate, br, zstd",
        },
    },
    # =========================
    # EDGE
    # =========================
    {
        "impersonate": "edge99",
        "headers": {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36 Edg/99.0.1150.30",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.9",
            "Accept-Encoding": "gzip, deflate, br",
            "Upgrade-Insecure-Requests": "1",
        },
    },
    {
        "impersonate": "edge101",
        "headers": {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36 Edg/101.0.1210.53",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.9",
            "Accept-Encoding": "gzip, deflate, br",
            "Upgrade-Insecure-Requests": "1",
        },
    },
    # =========================
    # SAFARI (MACOS)
    # =========================
    *[
        {
            "impersonate": name,
            "headers": {
                "User-Agent": ua,
                "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
                "Accept-Language": "en-us",
                "Accept-Encoding": "gzip, deflate, br",
            },
        }
        for name, ua in [
            (
                "safari153",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15",
            ),
            (
                "safari155",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Safari/605.1.15",
            ),
            (
                "safari170",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
            ),
            (
                "safari180",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15",
            ),
            (
                "safari184",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.4 Safari/605.1.15",
            ),
        ]
    ],
    # =========================
    # SAFARI (iOS)
    # =========================
    {
        "impersonate": "safari172_ios",
        "headers": {
            "User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
            "Accept-Language": "en-us",
            "Accept-Encoding": "gzip, deflate, br",
        },
    },
    {
        "impersonate": "safari180_ios",
        "headers": {
            "User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
            "Accept-Language": "en-us",
            "Accept-Encoding": "gzip, deflate, br",
        },
    },
    {
        "impersonate": "safari184_ios",
        "headers": {
            "User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 18_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.4 Mobile/15E148 Safari/604.1",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
            "Accept-Language": "en-us",
            "Accept-Encoding": "gzip, deflate, br",
        },
    },
]

for prof in profiles:
    impersonate = prof.get("impersonate")
    print(f"using {impersonate}")
    import random
    time_to_sleep = random.randint(1, 10)
    import time
    print(f"sleeping {time_to_sleep}...")
    time.sleep(time_to_sleep)
    headers = prof.get("headers")
    url = "https://www.linkedin.com/jobs/search/?currentJobId=4347732846&distance=25.0&geoId=105646813&keywords=%22scraping%22&origin=HISTORY"
    response = requests.get(url, headers=headers, impersonate=impersonate)
    soup = BeautifulSoup(response.text, "html.parser")
    job_list = soup.select("ul.jobs-search__results-list > li")
    print(f"Found {len(job_list)} job postings on the page.")
    
    if len(job_list) > 7:
        print("Success!")


# firefox 133, chrome99, chrome99_android