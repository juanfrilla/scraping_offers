package db

import (
	"database/sql"
	"fmt"
	"strings"

	"scraping_offers/go/models"
)

func SaveJobs(database *sql.DB, jobs []models.ScrapedJob) error {
	if len(jobs) == 0 {
		return nil
	}

	query := `
		INSERT INTO jobs (
			title,
			company,
			location,
			url,
			date_posted,
			modality,
			platform,
			keyword_appeared,
			logo_url
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)

		ON CONFLICT(platform, url) DO UPDATE SET
			title = excluded.title,
			company = excluded.company,
			location = excluded.location,
			date_posted = excluded.date_posted,
			modality = excluded.modality,
			keyword_appeared = excluded.keyword_appeared,
			logo_url = excluded.logo_url;
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

		_, err := stmt.Exec(
			job.Title,
			job.Company,
			job.Location,
			job.URL,
			job.DatePosted,
			job.Modality,
			job.Platform,
			job.KeywordAppeared,
			job.LogoURL,
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

func GetJobs(database *sql.DB) ([]models.ScrapedJob, error) {
	query := `
		SELECT
			title,
			company,
			location,
			url,
			date_posted,
			modality,
			platform,
			keyword_appeared,
			logo_url
		FROM jobs
		ORDER BY date_posted DESC;
	`

	rows, err := database.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.ScrapedJob

	for rows.Next() {
		var job models.ScrapedJob

		err := rows.Scan(
			&job.Title,
			&job.Company,
			&job.Location,
			&job.URL,
			&job.DatePosted,
			&job.Modality,
			&job.Platform,
			&job.KeywordAppeared,
			&job.LogoURL,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning job: %w", err)
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating jobs: %w", err)
	}

	return jobs, nil
}
