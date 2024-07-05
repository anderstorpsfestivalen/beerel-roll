package main

import (
	"fmt"
	"os"

	"github.com/wbergg/beerel-roll/dataimport"
	"github.com/wbergg/beerel-roll/db"
	"github.com/wbergg/beerel-roll/web"
)

func main() {
	fmt.Println("Starting beerel-roll")

	dbPath := "./db/db.sql"

	if os.Getenv("DATABASE_PATH") != "" {
		 fmt.Println("Using DATABASE_PATH:", os.Getenv("DATABASE_PATH"))
		 dbPath = os.Getenv("DATABASE_PATH") + "db.sql"
	}


	// Setup db
	d := db.Open(dbPath)

	// Check if DB is set up, if not, set it up (first time only)
	if d.Setup == 0 {
		fmt.Println("Looks like it's the first time - Populating DB...")
		err := dataimport.DbSetup(&d)
		if err == nil {
			fmt.Println("DB population sucess! Please rerun the program!")
		}
		os.Exit(0)
	}

	// Start webserver
	web.Start(&d)
}
