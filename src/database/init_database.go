package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

// Create target connection for the database
func GetConn() *pgx.Conn {
	connConfig := pgx.ConnConfig{
		Host:     "34.170.5.185",
		Port:     5432,
		Database: "rapidtransfer",
		User:     "postgres",
		Password: "postgres",
	}
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		fmt.Print("Failed to connect")
	}
	return conn
}

// Creates a public connection that all functions can use
var conn *pgx.Conn = GetConn()

// Executes returns str of .sql file when given the path
func execSQLFile(filename string) error {
	// Read the content of the SQL file
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Execute the SQL queries
	_, err = conn.Exec(string(content))
	if err != nil {
		return err
	}

}

// Inits all of the tables for the database
func InitializeDatabase() {
	const path = "database.sql"

	err := execSQLFile(path)

	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}

}
