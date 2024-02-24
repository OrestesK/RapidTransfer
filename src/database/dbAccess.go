package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

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
	conn.Exec("DROP TABLE transfer; DROP TABLE users; DROP TABLE friends; DROP TABLE user;")

	conn.Exec("CREATE TABLE IF NOT EXISTS transfer (userFrom INT NOT NULL, userTo INT NOT NULL, keyword VARCHAR(100))")
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT '', keyword VARCHAR(100), macaddr VARCHAR(100))")
	fmt.Printf("err: %v\n", err)
	conn.Exec("CREATE TABLE IF NOT EXISTS friends (orig_user INT NOT NULL, friend_id INT NOT NULL, total_transfers INT NOT NULL DEFAULT 0)")

	// _, err2 := conn.Exec("INSERT INTO users (name, keyword, macaddr) VALUES ($1, $2, $3)", "Emmet", "Keyword", "Something")
	// fmt.Printf("err2: %v\n", err2)
}

// func DropTables() {
// 	var conn *pgx.Conn = GetConn()
// 	conn.Exec("DROP TABLE transfer; DROP TABLE user; DROP TABLE friends;")
// }

func GenerateFriendCode() string {
	rand.Seed(time.Now().UnixNano())

	// Define characters that can be part of the friend code
	// You can customize this based on your requirements
	allowedChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Generate a random 8-character string
	var result string
	for i := 0; i < 8; i++ {
		result += string(allowedChars[rand.Intn(len(allowedChars))])
	}

	return result
}

func HandleAccountStartup() {
	conn := GetConn()
	macAddr, _ := getMacAddress()
	row := conn.QueryRow("SELECT * FROM users WHERE macaddr = $1", macAddr)

	var id int
	var name, keyword string
	var macaddr string
	err := row.Scan(&id, &name, &keyword, &macaddr)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Print("Error no row found for mac address!")
		}
	}

	if name == "" {
		// not found or empty username.
		// Prompt user for inputs
		fmt.Print("You appear to be a new user! ")
		fmt.Print("Enter your username to get started: ")
		_, err := fmt.Scan(&name)
		if err != nil {
			// fmt.Println("Error reading input:", err)
			os.Exit(1)
		}

		// fmt.Printf("Macddr: %s\n", macAddr)

		CreateAccount(name, macAddr)
		fmt.Print("Creating account!\n\n")

	} else {
		fmt.Print("You already exist! Logging you in!")
	}
}

func getMacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// Iterate through the network interfaces to find the MAC address
	for _, intf := range interfaces {
		if intf.HardwareAddr != nil {
			return intf.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("MAC address not found")
}

func CreateAccount(username string, macAddress string) {
	conn := GetConn()
	code := GenerateFriendCode()
	tag, err := conn.Exec("INSERT INTO users (name, keyword, macaddr) VALUES ($1, $2, $3)", username, code, macAddress)
	fmt.Printf("tag.RowsAffected(): %v\n", tag.RowsAffected())
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
}

// Retrieves a user's freind code based on their name, which is passed in
func GetUserFriendCode(name string) (userKey int) {
	conn := GetConn()
	err := conn.QueryRow("SELECT friendCode FROM users WHERE name=$1", (name)).Scan(&userKey)
	if err != nil {
		fmt.Println("Query failed")
	}
	return
}

// Retrieves a user's id based on their name, which is passed in
func GetUserID(name string) (userID int) {
	conn := GetConn()
	err := conn.QueryRow("SELECT id FROM users WHERE name=$1", (name)).Scan(&userID)
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

// Adds a friend to a senders friends list based on their friend code
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

func IsFriend(userName1, userName2 string) (isFriend bool) {
	isFriend = false
	conn := GetConn()
	user1Id := GetUserID(userName1)
	user2Id := GetUserID(userName2)
	err1 := conn.QueryRow("SELECT id FROM friends WHERE user_from=%s, user_to=%s", user1Id, user2Id).Scan(&isFriend)
	if err1 != nil {
		isFriend = true
		return
	}
	return
}

// Determines if two friends are mutual friends
func AreMutualFriends(userName1, userName2 string) (areMutuals bool) {
	areMutuals = false
	friend1 := IsFriend(userName1, userName2)
	friend2 := IsFriend(userName2, userName1)
	if friend1 == friend2 {
		areMutuals = true
		return
	}
	return
}

// Allows two users to send files to eachother
func PerformTransaction(senderName, recieverName string) (keyword string) {
	conn := GetConn()
	FromUserID := GetUserID(senderName)
	ToUserID := GetUserID(recieverName)
	if AreMutualFriends(senderName, recieverName) {
		keyword = GenerateFriendCode()
		err := conn.QueryRow("INSERT INTO transfer VALUES (%s,%s,%s)", FromUserID, ToUserID, keyword)
		if err != nil {
			fmt.Println("Query failed")
		}
	}
	return
}
