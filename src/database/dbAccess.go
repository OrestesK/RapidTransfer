package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

func GetConn() *pgx.Conn {
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

func InitializeDatabase() {

	var conn *pgx.Conn = GetConn()

	query, _ := ReadSQLFile("tansferDB.sql")

	conn.Query(query)

}

func ReadSQLFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading SQL file: %w", err)
	}
	return string(data), nil
}

// Retrieves a user's freind code based on their name, which is passed in
func GetUserKey(name string) int {
	conn := GetConn()
	var userKey int
	err := conn.QueryRow("SELECT friendCode FROM user WHERE name=%s", (name)).Scan(&userKey)
	if err != nil {
		fmt.Println("Query failed")
	}
	return userKey
}

// Retrieves a user's id based on their name, which is passed in
func GetUserID(name string) int {
	conn := GetConn()
	var userID int
	err := conn.QueryRow("SELECT id FROM user WHERE name=%s", (name)).Scan(&userID)
	if err != nil {
		fmt.Println("Query failed")
	}
	return userID
}
