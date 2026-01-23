# 🕵️‍♂️ Job Scraping Aggregator: Spain

A specialized web scraping tool designed to aggregate job postings from **LinkedIn**, **InfoJobs**, **Tecnoempleo**, **SimplyHired**, **RemoteRocketShip** and **RemoteOK**. This project automates the job seeking journey by centralizing offers into a single, searchable dashboard.

🚀 **Live Demo:** [scrapingoffersspain.streamlit.app](https://scrapingoffersspain.streamlit.app/)

---

## ⚠️ Residential Proxy Requirements on SimplyHired and RemoteRocketShip

- **SimplyHired**: Currently limited in the cloud demo. Due to strict security measures, it functions best in local environments or in production with residential proxies.
- **RemoteRocketShip**: The scraper has been specifically adapted to extract and display key job data directly within this app in production, reducing the need for residential proxies, but I can't access some information without residential proxies, specifically the job offer url.

## ✨ Features

- **Multi-Source Aggregation:** Scrapes four major job boards focused on the Spanish market.
- **Interactive Dashboard:** Built with Streamlit for real-time filtering and data visualization.
- **Fast Performance:** Optimized dependency management and scraping logic.

## 🛠️ Installation & Local Setup

To run this project locally, ensure you have Python installed. This project uses `uv` for high-performance package management.

### 1. Clone the Repository

```bash
git clone https://github.com/juanfrilla/scraping_offers.git
cd scraping_offers
```

### 2. Set up environment variables

```bash
echo "ENV=local" > .env
```

### 3. Create and activate a virtual environment Using uv:

```bash
uv venv .venv
# On macOS/Linux:
source .venv/bin/activate
# On Windows:
# .venv\Scripts\activate
```

### 4. Install dependencies

```bash
uv pip install -r requirements.txt
```

### 5. Run the app

```bash
streamlit run main.py
```
