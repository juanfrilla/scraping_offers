import sqlite3


def get_jobs(database_file):

    with sqlite3.connect(database_file, timeout=10) as conn:
        conn.row_factory = sqlite3.Row

        cursor = conn.execute("""
            SELECT
                id,
                title,
                company,
                location,
                url,
                date_posted,
                modality,
                platform,
                keyword_appeared,
                logo_url,
                external_id,
                created_at,
                updated_at
            FROM jobs
            ORDER BY date_posted DESC
        """)

        return [dict(row) for row in cursor.fetchall()]
