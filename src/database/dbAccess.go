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

var currentUser *User

type User struct {
	id      int
	name    string
	keyword string
	macaddr string
}

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
	fmt.Print("Hello! ")
	var conn *pgx.Conn = GetConn()
	// conn.Exec("DROP TABLE transfer; DROP TABLE users; DROP TABLE friends; DROP TABLE user;")

	conn.Exec("CREATE TABLE IF NOT EXISTS transfer (id SERIAL PRIMARY KEY, userFrom INT NOT NULL, userTo INT NOT NULL, keyword VARCHAR(100), address VARCHAR(100), filename VARCHAR(100))")
	conn.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL UNIQUE, keyword VARCHAR(100), macaddr VARCHAR(100))")
	conn.Exec("CREATE TABLE IF NOT EXISTS friends (orig_user INT NOT NULL, friend_id INT NOT NULL, total_transfers INT NOT NULL DEFAULT 0)")

}

func generateFriendCode() string {
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

func GetUserDetails() (int, string, string, string) {
	if currentUser != nil {
		return currentUser.id, currentUser.name, currentUser.keyword, currentUser.macaddr
	}
	panic("WTF")
}

func GetPendingTransfers() {
	conn := GetConn()
	rows, _ := conn.Query("SELECT users.name as host, transfer.keyword, filename FROM transfer INNER JOIN users ON users.id = transfer.userFrom WHERE userTo = $1", 3)
	for rows.Next() {

		var host, keyword, filename string

		err := rows.Scan(&host, &keyword, &filename)
		if err != nil {
			// Handle the error
			fmt.Println("Error scanning row:", err)
			continue
		}

		// Print the values or perform any other desired operation
		fmt.Printf("Host: %s, Keyword: %s, Filename: %s\n", host, keyword, filename)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}
}

/*
*

	Called on process start. This will find the user if it exists in the database. If not, it will ask to create an account.
*/
func HandleAccountStartup() {
	conn := GetConn()
	macAddr, erro := getMacAddress()
	if erro != nil {
		panic(erro)
	}
	row := conn.QueryRow("SELECT * FROM users WHERE macaddr = $1", macAddr)

	var id int
	var name, keyword string
	var macaddr string
	err := row.Scan(&id, &name, &keyword, &macaddr)

	//
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Print("Error no row found for mac address!")
		}
	}

	if name == "" {
		// not found or empty username.
		// Prompt user for inputs
		fmt.Print("You appear to be a new user!\nEnter your username to get started: ")
		_, err := fmt.Scan(&name)
		if err != nil {
			// fmt.Println("Error reading input:", err)
			os.Exit(1)
		}

		// fmt.Printf("Macddr: %s\n", macAddr)

		CreateAccount(name, macAddr)
		fmt.Print("Creating account!\n\n")
		currentUser = &User{
			id:      1,
			name:    name,
			keyword: keyword,
			macaddr: macaddr,
		}

	} else {
		fmt.Printf("You already exist! Logging you in as %s!\n\n", name)
		currentUser = &User{
			id:      1,
			name:    name,
			keyword: keyword,
			macaddr: macaddr,
		}
	}
}

func GetCurrentUser() User {
	if currentUser != nil {
		return *currentUser
	}
	return User{}
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

/*
*

	Creates an account in the database with specifid username/macaddr.
	MacAddr is unique to the computer, and is used on startup to indentify the pc.
*/
func CreateAccount(username string, macAddress string) {
	conn := GetConn()
	code := generateFriendCode()
	_, err := conn.Exec("INSERT INTO users (name, keyword, macaddr) VALUES ($1, $2, $3)", username, code, macAddress)
	if err != nil {
		fmt.Printf("Failed at CreateAccount: %s", err)
	}
}

// Retrieves a user's freind code based on their name, which is passed in
func GetUserFriendCode(name string) (userKey string) {
	conn := GetConn()
	err := conn.QueryRow("SELECT keyword FROM users WHERE name=$1", name).Scan(&userKey)
	if err != nil {
		fmt.Print("Failed at GetUserFriendCode")
	}
	return
}

// Retrieves a user's id based on their name, which is passed in
func GetUserID(name string) (userID int) {
	conn := GetConn()
	fmt.Printf("Name: %s\n", name)
	err := conn.QueryRow("SELECT id FROM users WHERE name=$1", (name)).Scan(&userID)
	if err != nil {
		fmt.Print("Failed at GetUserID")
		panic(err)
	}
	return
}

// Retrieves a user's name based on their id, which is passed in
func GetUserNameByID(id int) (userName string) {
	conn := GetConn()
	err := conn.QueryRow("SELECT name FROM users WHERE id=$1", (id)).Scan(&userName)
	if err != nil {
		fmt.Print("Failed at GetUserNameByID")
		panic(err)
	}
	return
}

// Retrieves a user's name based on their friend code, which is passed in
func GetUserNameByFriendCode(friendCode int) (userName string) {
	conn := GetConn()
	err := conn.QueryRow("SELECT name FROM user WHERE id=$1", (friendCode)).Scan(&userName)
	if err != nil {
		fmt.Print("Failed at GetUserNameByFriendCode")
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
	err := conn.QueryRow("SELECT uidFrom, uidTo FROM transaction WHERE keyword=$1", (keyword)).Scan(&userFromID, &userToID)
	if err != nil {
		fmt.Print("Failed at GetTransaction")
	}
	userFromName := GetUserNameByID(userFromID)
	userToName := GetUserNameByID(userToID)
	names = []string{userFromName, userToName}
	return
}

// Adds a friend to a senders friends list based on their friend code
func AddFriend(friendCode string, senderName string) (success bool) {
	conn := GetConn()
	var friendID int
	fmt.Printf("friendcode: %s\n", friendCode)
	err := conn.QueryRow("SELECT id FROM users WHERE keyword=$1;", friendCode).Scan(&friendID)
	fmt.Printf("FriendID From Code: %d\n", friendID)
	if err != nil {
		fmt.Print("Failed at AddFriend")
		return false
	}
	senderID := GetUserID(senderName)

	row := conn.QueryRow("SELECT COUNT(*) FROM friends WHERE orig_user = $1 AND friend_id = $2", senderID, friendID)
	var count int
	row.Scan(&count)
	if count > 0 {
		fmt.Print("Friend is already added!")
		return
	}

	_, err2 := conn.Exec("INSERT INTO friends VALUES ($1, $2)", senderID, friendID)
	_, err3 := conn.Exec("INSERT INTO friends VALUES ($1, $2)", friendID, senderID)
	if err2 != nil || err3 != nil {
		fmt.Print("Failed at AddFriend 2")
	}
	return true
}

// Deletes a one way friendship between two users
func DeleteFriend(senderName string, recieverName string) (deletedFriend string) {
	conn := GetConn()
	senderId := GetUserID(senderName)
	recieverId := GetUserID(recieverName)
	err := conn.QueryRow("DELETE FROM friends WHERE user_to=$1, user_from=$1", senderId, recieverId).Scan(&deletedFriend)
	if err != nil {
		fmt.Println("Query failed")
	}
	return
}

func IsFriend(userName1 string, userName2 string) (isFriend bool) {
	isFriend = false
	conn := GetConn()
	user1Id := GetUserID(userName1)
	user2Id := GetUserID(userName2)
	err1 := conn.QueryRow("SELECT id FROM friends WHERE user_from=$1, user_to=$2", user1Id, user2Id).Scan(&isFriend)
	if err1 != nil {
		isFriend = true
		return
	}
	return
}

func GetAddressFromTransactionPhrase(phrase string) string {
	conn := GetConn()
	row := conn.QueryRow("SELECT address FROM transfer WHERE keyword = $1;", phrase)
	var address string
	err := row.Scan(&address)
	if err != nil {
		panic(err)
	}
	return address
}

func GetFileNameFromTransactionPhrase(phrase string) string {
	conn := GetConn()
	row := conn.QueryRow("SELECT filename FROM transfer WHERE keyword = $1;", phrase)
	var filename string
	err := row.Scan(&filename)
	if err != nil {
		panic(err)
	}
	return filename
}

func DeleteTransactionWithAddress(address string) {
	conn := GetConn()
	_, err := conn.Exec("DELETE FROM transfer WHERE address = $1;", address)
	if err != nil {
		panic(err)
	}
}

// Determines if two friends are mutual friends
func AreMutualFriends(userName1 string, userName2 string) (areMutuals bool) {
	areMutuals = false
	return IsFriend(userName1, userName2) && IsFriend(userName2, userName1)
}

// Allows two users to send files to eachother
func PerformTransaction(senderName string, recieverName string, address string, filename string) (keyword string) {
	conn := GetConn()
	FromUserID := GetUserID(senderName)
	ToUserID := GetUserID(recieverName)

	if AreMutualFriends(senderName, recieverName) {
		phrase := generateFriendCode()
		_, err := conn.Exec("INSERT INTO transfer (userFrom, userTo, keyword, address, filename) VALUES ($1,$2,$3,$4,$5)", FromUserID, ToUserID, phrase, address, filename)
		if err != nil {
			fmt.Print("Failed at PerformTransaction")
			panic(err)
		}
		return phrase
	}
	fmt.Print("Failed to create transaction")
	return ""
}

func GetFriendsList(username string) (friendsList []string) {
	conn := GetConn()
	userId := GetUserID(username)

	conn.QueryRow(`
	
	SELECT ARRAY_AGG(users.name)
	FROM users 
	
	INNER JOIN friends ON friends.friend_id=users.id
	WHERE orig_user=$1 
	
	GROUP BY friends.orig_user`, userId).Scan(&friendsList)
	fmt.Println(friendsList)
	return
}
