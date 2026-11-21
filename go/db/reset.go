package db

import (
	"fmt"
	"os"
	"scraping_offers/go/constants"
)

func ResetDB() error {
	// Remove existing database file
	if err := os.Remove(constants.DataDB); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing database file: %w", err)
	}

	fmt.Println("Database file removed:", constants.DataDB)

	// Create a fresh database
	database, err := InitDB(constants.DataDB)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}
	defer database.Close()

	fmt.Println("Database recreated:", constants.DataDB)

	// Check indexes
	rows, err := database.Query("PRAGMA index_list('jobs')")
	if err != nil {
		return fmt.Errorf("checking indexes: %w", err)
	}
	defer rows.Close()

	fmt.Println("Indexes in jobs table:")

	for rows.Next() {
		var seq int
		var name string
		var unique int
		var origin string
		var partial int

		if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			return fmt.Errorf("scanning indexes: %w", err)
		}

		fmt.Printf(
			"Name: %s | Unique: %d | Origin: %s | Partial: %d\n",
			name,
			unique,
			origin,
			partial,
		)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterating indexes: %w", err)
	}

	return nil
}
