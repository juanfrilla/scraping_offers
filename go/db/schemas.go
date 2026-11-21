package db

import "database/sql"

func CreateSchema(database *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS jobs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL DEFAULT '',
            company TEXT NOT NULL DEFAULT '',
            location TEXT NOT NULL DEFAULT '',
            url TEXT NOT NULL,
            date_posted INTEGER,
            modality TEXT NOT NULL DEFAULT '',
            platform TEXT NOT NULL,
            keyword_appeared TEXT NOT NULL DEFAULT '',
            logo_url TEXT NOT NULL DEFAULT '',
            external_id TEXT NOT NULL,

            created_at INTEGER NOT NULL DEFAULT (unixepoch()),
            updated_at INTEGER NOT NULL DEFAULT (unixepoch()),

            UNIQUE(platform, external_id)
        );

        CREATE INDEX IF NOT EXISTS idx_jobs_date_posted
        ON jobs(date_posted);
    `

	_, err := database.Exec(query)
	return err
}
