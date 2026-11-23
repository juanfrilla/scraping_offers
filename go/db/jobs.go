package db

import (
	"database/sql"
	"fmt"
	"scraping_offers/go/models"
	"strings"
	"time"
)

func SaveJobs(database *sql.DB, jobs []models.ScrapedJob) error {
	if len(jobs) == 0 {
		return nil
	}

	query := `
INSERT INTO jobs (
    title, company, location, url, date_posted,
    modality, platform, keyword_appeared, logo_url, external_id,
    updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, unixepoch())
ON CONFLICT(platform, external_id) DO UPDATE SET
    title            = excluded.title,
    company          = excluded.company,
    location         = excluded.location,
    url              = excluded.url,
    date_posted      = excluded.date_posted,
    modality         = excluded.modality,
    keyword_appeared = excluded.keyword_appeared,
    logo_url         = excluded.logo_url,
    updated_at       = unixepoch();
`

	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	for i, job := range jobs {
		if strings.TrimSpace(job.URL) == "" {
			return fmt.Errorf("job at index %d has an empty URL", i)
		}

		if strings.TrimSpace(job.Platform) == "" {
			return fmt.Errorf("job at index %d has an empty platform", i)
		}

		if strings.TrimSpace(job.ExternalID) == "" {
			return fmt.Errorf("job at index %d has an empty external ID", i)
		}

		var datePosted any

		if !job.DatePosted.IsZero() {
			datePosted = job.DatePosted.UTC().Unix()
		}

		_, err := stmt.Exec(
			job.Title,
			job.Company,
			job.Location,
			job.URL,
			datePosted,
			job.Modality,
			job.Platform,
			job.KeywordAppeared,
			job.LogoURL,
			job.ExternalID,
		)

		if err != nil {
			return fmt.Errorf("saving job %s: %w", job.URL, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func DeleteJobsNotUpdatedToday(database *sql.DB, now time.Time, loc *time.Location) (int64, error) {
	now = now.In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	startOfNextDay := startOfDay.AddDate(0, 0, 1)

	const query = `
        DELETE FROM jobs
        WHERE updated_at IS NULL
           OR updated_at <  ?
           OR updated_at >= ?
    `
	res, err := database.Exec(query, startOfDay.Unix(), startOfNextDay.Unix())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetJobsCount(database *sql.DB) (int, error) {
	var count int

	err := database.QueryRow("SELECT COUNT(*) FROM jobs").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
