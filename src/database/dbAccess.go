package database

import (
	"fmt"
	// "os"

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
func GetTransfer(keyword string) (names []string) {
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
