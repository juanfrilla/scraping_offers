package main

import (
	"fmt"
	"os"

	"scraping_offers/go/constants"
	"scraping_offers/go/db"
)

func main() {

	if err := constants.InitPaths(); err != nil {
		fmt.Println("Error initializing paths:", err)
		os.Exit(1)
	}
	if err := db.ResetDB(); err != nil {
		fmt.Println("Error resetting database:", err)
		os.Exit(1)
	}

	fmt.Println("Database successfully reset!")
}
