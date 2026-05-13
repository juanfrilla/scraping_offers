# 🕵️‍♂️ Job Scraping Aggregator: Spain

A specialized web scraping tool designed to aggregate job postings from **LinkedIn**, **InfoJobs**, **Tecnoempleo**, **SimplyHired**, **RemoteRocketShip**, and **RemoteOK**. This project automates the job seeking journey by centralizing offers into a single, searchable dashboard.

🚀 **Live Demo:** [scrapingoffersspain.streamlit.app](https://scrapingoffersspain.streamlit.app/)

---

## ⚠️ Residential Proxy Requirements

- **SimplyHired:** Currently limited in the cloud demo. Due to strict security measures, it functions best in local environments or in production with residential proxies.
- **RemoteRocketShip:** The scraper has been specifically adapted to extract and display key job data directly within this app in production, reducing the need for residential proxies — however, some information (specifically the job offer URL) remains inaccessible without residential proxies.

---

## ✨ Features

- **Multi-Source Aggregation:** Scrapes major job boards focused on the Spanish market.
- **Interactive Dashboard:** Built with Streamlit for real-time filtering and data visualization.
- **Fast Performance:** Scrapers migrated to **Go (Golang)** for significantly improved speed and efficiency.

---

## 🏗️ Architecture

The project can be used in two modes:

| Mode | Command | Description |
|------|---------|-------------|
| Full stack | `streamlit run main.py` | Launches the dashboard and spawns the Go scraper in the background |
| Scraper only | `go run main.go` | Runs scrapers standalone and saves results to `data/jobs.json` |

---

## 🛠️ Installation & Local Setup

### Prerequisites

- [Python](https://www.python.org/downloads/) (for the Streamlit dashboard)
- [Go](https://go.dev/dl/) 1.21+ (for the scrapers)
- [`uv`](https://github.com/astral-sh/uv) (high-performance Python package manager)

---

### 1. Clone the Repository

```bash
git clone https://github.com/juanfrilla/scraping_offers.git
cd scraping_offers
```

### 2. Set Up Environment Variables

```bash
echo "ENV=local" > .env
```

### 3. Create and Activate a Virtual Environment

```bash
uv venv .venv

# macOS/Linux:
source .venv/bin/activate

# Windows:
# .venv\Scripts\activate
```

### 4. Install Python Dependencies

```bash
uv pip install -r requirements.txt
```

### 5. Install Go Dependencies

Dependencies are managed automatically via Go Modules. To download them explicitly before running:

```bash
go mod download
```

> On `go run` or `go build`, Go will fetch any missing dependencies automatically.

---

## 🚀 Running the App

### Option A — Full Stack (Dashboard + Scrapers)

```bash
streamlit run main.py
```

The Streamlit app automatically launches the Go scraper in the background. Results are fetched and displayed in real time in the dashboard.

---

### Option B — Scraper Only (Go)

Run the scrapers independently without the dashboard:

```bash
go run main.go
```

Or compile and run the binary:

```bash
# Build
go build -o scraper main.go

# Run
./scraper        # macOS/Linux
# scraper.exe    # Windows
```

Scraped jobs are saved to `data/jobs.json`