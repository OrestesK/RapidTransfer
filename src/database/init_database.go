package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx"
)

var conn *pgx.Conn
var connMutex sync.Mutex

// Create target connection for the database
func GetConn() (*pgx.Conn, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if conn != nil {
		return conn, nil
	}

	connConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "rapidtransfer",
		User:     "swen344",
		Password: "Forzano17**",
	}

	newConn, err := pgx.Connect(connConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %v", err)
	}

	conn = newConn
	return conn, nil
}

// Executes returns str of .sql file when given the path
func execSQLFile(filename string, conn *pgx.Conn) error {
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

	return nil
}

// Inits all of the tables for the database
func InitializeDatabase() {
	const path = "src/database/database.sql"

	// Explicitly initialize the connection
	conn, err := GetConn()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}

	// Execute the SQL file
	err = execSQLFile(path, conn)
	if err != nil {
		fmt.Println("Error executing SQL file:", err)
		os.Exit(1)
	}
}
