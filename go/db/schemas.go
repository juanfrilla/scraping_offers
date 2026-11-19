package db

import "database/sql"

func CreateSchema(database *sql.DB) error {

	//TODO mirar esto detenidamente
	// REVISAR SI EL NOT NULL DEL EXT ID SE PUEDE PONER ASI
	query := `
        CREATE TABLE IF NOT EXISTS jobs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL DEFAULT '',
            company TEXT NOT NULL DEFAULT '',
            location TEXT NOT NULL DEFAULT '',
            url TEXT NOT NULL,
            date_posted TEXT,
            modality TEXT NOT NULL DEFAULT '',
            platform TEXT NOT NULL,
            keyword_appeared TEXT NOT NULL DEFAULT '',
            logo_url TEXT NOT NULL DEFAULT '',
            external_id TEXT NOT NULL, 

            created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

            UNIQUE(platform, url)
        );

        CREATE INDEX IF NOT EXISTS idx_jobs_platform
        ON jobs(platform);

        CREATE INDEX IF NOT EXISTS idx_jobs_date_posted
        ON jobs(date_posted);

        CREATE INDEX IF NOT EXISTS idx_jobs_external_id
        ON jobs(platform, external_id);
    `

	_, err := database.Exec(query)
	return err
}
