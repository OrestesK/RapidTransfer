package database

import (
	"fmt"

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

// Retrieves a user's name based on their id, which is passed in
func GetUserNameByID(id int) (userName string) {
	conn := GetConn()
	err := conn.QueryRow("SELECT name FROM user WHERE id=%s", (id)).Scan(&userName)
	if err != nil {
		fmt.Println("Query failed")
	}
	return
}

// Retrieves a user's name based on their friend code, which is passed in
func GetUserNameByFriendCode(friendCode int) (userName string) {
	conn := GetConn()
	err := conn.QueryRow("SELECT name FROM user WHERE id=%s", (friendCode)).Scan(&userName)
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
	userFromName := GetUserNameByID(userFromID)
	userToName := GetUserNameByID(userToID)
	names = []string{userFromName, userToName}
	return
}

func AddFriend(friendCode int, senderName string) (friendName string) {
	conn := GetConn()
	var friendID int
	err := conn.QueryRow("SELECT id FROM user WHERE friendCode=%s", (friendCode)).Scan(&friendID)
	if err != nil {
		fmt.Println("Query failed")
	}
	senderID := GetUserID(senderName)
	err2 := conn.QueryRow("INSERT INTO friends VALUES (%s,%s)", senderID, friendID)
	if err2 != nil {
		fmt.Println("Query failed")
	}
	friendName = GetUserNameByID(friendID)
	return
}

func AreMutualFriends(userName1, userName2 string) (areMutuals bool) {
	conn := GetConn()
	user1Id := GetUserID(userName1)
	user2Id := GetUserID(userName2)
	err1 := conn.QueryRow("SELECT id FROM friends WHERE user_from=%s, user_to=%s", user1Id, user2Id).Scan(&areMutuals)
	if err1 != nil {
		areMutuals = false
		return
	}
	err2 := conn.QueryRow("SELECT id FROM friends WHERE user_from=%s, user_to=%s", user2Id, user1Id).Scan(&areMutuals)
	if err2 != nil {
		areMutuals = false
		return
	}
	areMutuals = true
	return
}
