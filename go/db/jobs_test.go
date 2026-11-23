package db

import (
	"database/sql"
	"scraping_offers/go/models"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("opening test database: %v", err)
	}

	if err := CreateSchema(database); err != nil {
		database.Close()
		t.Fatalf("creating schema: %v", err)
	}

	return database
}

func TestSaveJobs(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	jobs := []models.ScrapedJob{
		{
			Title:      "Python Developer",
			Company:    "Acme",
			Location:   "Remote",
			URL:        "https://example.com/job/1",
			Platform:   "LinkedIn",
			ExternalID: "1",
		},
		{
			Title:      "Go Developer",
			Company:    "Acme",
			Location:   "Remote",
			URL:        "https://example.com/job/2",
			Platform:   "LinkedIn",
			ExternalID: "2",
		},
	}

	err := SaveJobs(database, jobs)
	if err != nil {
		t.Fatalf("SaveJobs() error = %v", err)
	}

	count, err := GetJobsCount(database)
	if err != nil {
		t.Fatalf("GetJobsCount() error = %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 jobs, got %d", count)
	}
}

func TestSaveJobsUpsert(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	job := models.ScrapedJob{
		Title:      "Python Developer",
		Company:    "Old Company",
		Location:   "Remote",
		URL:        "https://example.com/job/1",
		Platform:   "LinkedIn",
		ExternalID: "123",
	}

	err := SaveJobs(database, []models.ScrapedJob{job})
	if err != nil {
		t.Fatalf("first SaveJobs() error = %v", err)
	}

	job.Company = "New Company"

	err = SaveJobs(database, []models.ScrapedJob{job})
	if err != nil {
		t.Fatalf("second SaveJobs() error = %v", err)
	}

	count, err := GetJobsCount(database)
	if err != nil {
		t.Fatalf("GetJobsCount() error = %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 job after upsert, got %d", count)
	}

	var company string

	err = database.QueryRow(`
        SELECT company
        FROM jobs
        WHERE platform = ? AND external_id = ?
    `, "LinkedIn", "123").Scan(&company)

	if err != nil {
		t.Fatalf("querying job: %v", err)
	}

	if company != "New Company" {
		t.Errorf("expected company %q, got %q", "New Company", company)
	}
}

func TestDeleteJobsNotUpdatedToday(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	loc, err := time.LoadLocation("Atlantic/Canary")
	if err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 11, 23, 12, 0, 0, 0, loc)

	yesterday := now.AddDate(0, 0, -1).Unix()
	today := now.Unix()

	_, err = database.Exec(`
        INSERT INTO jobs (
            title, company, location, url,
            platform, external_id, updated_at
        )
        VALUES
            ('Old job', 'Company', 'Remote', 'https://old.com',
             'test', 'old', ?),
            ('Current job', 'Company', 'Remote', 'https://current.com',
             'test', 'current', ?)
    `, yesterday, today)

	if err != nil {
		t.Fatalf("inserting test jobs: %v", err)
	}

	deleted, err := DeleteJobsNotUpdatedToday(database, now, loc)
	if err != nil {
		t.Fatalf("DeleteJobsNotUpdatedToday() error = %v", err)
	}

	if deleted != 1 {
		t.Errorf("expected 1 deleted job, got %d", deleted)
	}

	count, err := GetJobsCount(database)
	if err != nil {
		t.Fatalf("GetJobsCount() error = %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 remaining job, got %d", count)
	}
}
