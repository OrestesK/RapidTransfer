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
func GetUserFriendCode(name string) (userKey int) {
	conn := GetConn()
	err := conn.QueryRow("SELECT friendCode FROM user WHERE name=%s", (name)).Scan(&userKey)
	if err != nil {
		fmt.Println("Query failed")
	}
	return
}

// Retrieves a user's id based on their name, which is passed in
func GetUserID(name string) (userID int) {
	conn := GetConn()
	err := conn.QueryRow("SELECT id FROM user WHERE name=%s", (name)).Scan(&userID)
	if err != nil {
		fmt.Println("Query failed")
	}
	return
}

// Retrieves a user's name based on their friend code, which is passed in
func GetUserName(id int) (userName string) {
	conn := GetConn()
	err := conn.QueryRow("SELECT name FROM user WHERE id=%s", (id)).Scan(&userName)
	if err != nil {
		fmt.Println("Query failed")
	}
	return
}

/*
Retrieves the names of two users who have had a transaction with eachother. This
function does this by reading the specific keyword associated with the transaction
*/
func GetTransaction(keyword string) (names []string) {
	conn := GetConn()
	var userFromID int
	var userToID int
	err := conn.QueryRow("SELECT uidFrom, uidTo FROM transaction WHERE keyword=%s", (keyword)).Scan(&userFromID, &userToID)
	if err != nil {
		fmt.Println("Query failed")
	}
	userFromName := GetUserName(userFromID)
	userToName := GetUserName(userToID)
	names = []string{userFromName, userToName}
	return
}
