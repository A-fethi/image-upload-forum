package database

import (
	"database/sql"
	"fmt"
	"forum/app/config"
	"os"
)

func CreateTables(db *sql.DB) {
	statement, err := os.ReadFile("./app/db/schema.sql")
	if err != nil {
		config.Logger.Println("Error Reading SQL Schema: ", err)
		os.Exit(1)
	}
	_, err = db.Exec(string(statement))
	if err != nil {
		config.Logger.Println("Error Executing SQL Schema To Create Tables: ", err)
		os.Exit(1)
	}

	fmt.Println("Database Connected successfully!")
}
