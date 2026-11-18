# 🕵️‍♂️ Job Scraping Aggregator: Spain

A specialized web scraping tool designed to aggregate job postings from **LinkedIn**, **InfoJobs**, **Tecnoempleo**, **SimplyHired**, **RemoteRocketShip**, **RemoteOK** and **Ashby**.

The project automates job discovery by collecting offers from multiple sources and centralizing them in a single, searchable dashboard.

🚀 **Live Demo:** [scrapingoffersspain.streamlit.app](https://scrapingoffersspain.streamlit.app/)

---

## ⚠️ Residential Proxy Requirements

- **SimplyHired:** Currently limited in the cloud demo. Due to strict anti-bot measures, it works best in local environments or production environments with residential proxies.
- **RemoteRocketShip:** The scraper has been adapted to extract and display key job information directly in the production app, reducing the need for residential proxies. However, some information, particularly job offer URLs, may remain inaccessible without them.

---

## ✨ Features

- **Multi-Source Aggregation:** Collects job postings from major job boards, with a focus on the Spanish market.
- **Interactive Dashboard:** Streamlit interface for browsing, filtering, and exploring job offers.
- **Go Scraping Engine:** Scrapers implemented in Go for improved performance and efficiency.
- **Manual Scraping:** Launch a new scraping run directly from the dashboard.
- **Automated Daily Scraping:** GitHub Actions runs the scrapers daily and updates the repository with the latest results.
- **CI/CD Pipeline:** Automatically builds the Go scraper and validates the Python application on every push to the main branch.
- **Logging and Error Handling:** Captures execution logs and handles scraping failures.

---

## 🏗️ Architecture

The project separates scraping, presentation, and automation responsibilities.

| Component | Technology | Responsibility |
|---|---|---|
| Scraping engine | Go | Extracts job postings from multiple sources |
| Dashboard | Python + Streamlit | Displays and filters collected offers |
| Data storage | JSON | Stores scraped job postings in `data/jobs.json` |
| CI/CD | GitHub Actions | Builds, validates, and automates scraping |

### Execution modes

| Mode | Command | Description |
|---|---|---|
| Dashboard | `streamlit run main.py` | Launches the Streamlit application |
| Scraper only | `go run ./go/main.go` | Runs the Go scrapers independently |
| Compiled scraper | `./go/bin/scraper_go` | Executes the compiled Go binary |

### Automated workflow

GitHub Actions handles three different scenarios:

1. **Push to `main`:** Builds the Go binary and validates the Python application. Scraping is not executed.
2. **Scheduled execution:** Runs the Go scraper daily at 00:00 UTC and commits updated data to the repository.
3. **Manual execution:** Allows triggering the complete scraping workflow from the GitHub Actions interface.

The Streamlit dashboard also provides a **Force Scraping** button to execute the Go scraper directly from the application environment.

---

## 🛠️ Installation & Local Setup

### Prerequisites

- [Python](https://www.python.org/downloads/)
- [Go](https://go.dev/dl/) 1.22+
- [`uv`](https://github.com/astral-sh/uv) (Python package manager)

### 1. Clone the repository

```bash
git clone https://github.com/juanfrilla/scraping_offers.git
cd scraping_offers
```

### 2. Configure environment variables

Create a `.env` file:

```env
ENV=local
```

### 3. Create and activate a virtual environment

```bash
uv venv .venv
```

**macOS/Linux:**

```bash
source .venv/bin/activate
```

**Windows (PowerShell):**

```powershell
.venv\Scripts\Activate.ps1
```

### 4. Install Python dependencies

```bash
uv pip install -r requirements.txt
```

### 5. Download Go dependencies

```bash
go mod download
```

Go Modules manages the project's Go dependencies. Missing dependencies can also be downloaded automatically during compilation.

---

## 🚀 Running the Application

### Option A — Streamlit Dashboard

```bash
streamlit run main.py
```

The dashboard loads the available job postings from `data/jobs.json`.

You can trigger a new scraping run using the **Force Scraping** button in the sidebar. This executes the Go scraper directly from the application environment and refreshes the dashboard when the process finishes.

### Option B — Standalone Go Scraper

Run the scraper independently, without launching Streamlit:

```bash
go run ./go/main.go
```

Alternatively, compile and execute the binary.

**Build:**

```bash
go build -o ./go/bin/scraper_go ./go/main.go
```

**Run (macOS/Linux):**

```bash
./go/bin/scraper_go
```

**Run (Windows):**

```powershell
go build -o ./go/bin/scraper_go.exe ./go/main.go
.\go\bin\scraper_go.exe
```

Scraped job postings are saved to:

```text
data/jobs.json
```

---

## ⚙️ GitHub Actions

The project uses GitHub Actions to automate builds, validation, and scraping.

### Daily scraping

The workflow runs automatically every day at **00:00 UTC**.

After a successful scraping run, updated job data and the compiled Go binary are committed to the repository.

### Manual workflow execution

You can also trigger the complete workflow whenever you want:

1. Open the repository on GitHub.
2. Navigate to **Actions**.
3. Select **Build Go Scraper & Validate Python**.
4. Click **Run workflow**.
5. Select the `main` branch and confirm.

### Continuous Integration

Every push to `main` triggers the build and validation steps without executing the scrapers.

This keeps regular development pushes lightweight while preserving the automated daily scraping process.

---

## 📁 Project Structure

```text
scraping_offers/
│
├── go/
│   ├── main.go
│   └── bin/
│       └── scraper_go
│
├── py/
│   ├── constants.py
│   ├── logger.py
│   └── utils.py
│
├── data/
│   └── jobs.json
│
├── logs/
│
├── main.py
├── requirements.txt
├── go.mod
├── go.sum
├── .env
└── .github/
    └── workflows/
        └── ...
```

---

## 📌 Notes

- The dashboard and scraping engine can be executed independently.
- GitHub Actions maintains the repository's daily dataset.
- Manual scraping from Streamlit executes in the application's environment; its filesystem changes are not automatically committed to GitHub.
- Scraping availability depends on each job board's anti-bot protections and access restrictions.