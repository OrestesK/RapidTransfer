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
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		fmt.Print("Failed to connect")
	}
	return conn
}

func InitializeDatabase() {
	fmt.Print("Hello")
	var conn *pgx.Conn = GetConn()
	conn.Exec("CREATE TABLE IF NOT EXISTS transfer (userFrom INT NOT NULL, userTo INT NOT NULL, keyword VARCHAR(100))")
	conn.Exec("CREATE TABLE IF NOT EXISTS user (int SERIAL INT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT '', keyword VARCHAR(100))")
	conn.Exec("CREATE TABLE IF NOT EXISTS friends (orig_user INT NOT NULL, friend_id INT NOT NULL, total_transfers INT NOT NULL DEFAULT 0)")

}

func DropTables() {
	var conn *pgx.Conn = GetConn()
	conn.Exec("DROP TABLE transfer; DROP TABLE user; DROP TABLE friende;")
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
