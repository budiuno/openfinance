package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	// Connection string for PostgreSQL
	connStr := "host=localhost port=5432 user=root password=root dbname=openfinanceDB sslmode=disable"

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return err
	}
	// defer db.Close()

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return err
	}

	fmt.Println("Successfully connected to PostgreSQL database!")
	DB = db

	return nil
}
