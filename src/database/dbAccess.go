package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

func getConn() *pgx.Conn {
	connConfig := pgx.ConnConfig{
		Host:     "34.170.5.185",
		Port:     5432,
		Database: "rapidtransfer",
		User:     "postgres",
		Password: "postgres",
	}
	conn, _ := pgx.Connect(connConfig)
	return conn
}

func initializeDatabase() {

	var conn *pgx.Conn = getConn()

	query, _ := readSQLFile("tansferDB.sql")

	conn.Query(query)

}

func readSQLFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading SQL file: %w", err)
	}
	return string(data), nil
}

