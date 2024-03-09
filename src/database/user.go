package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

var currentUser *User

type User struct {
	id      int
	name    string
	keyword string
	macaddr string
}

// Retrieves a user's name based on their id, which is passed in
func GetUserNameByID(id int) (userName string) {
	err := conn.QueryRow("SELECT name FROM users WHERE id=$1", id).Scan(&userName)
	if err != nil {
		fmt.Print("Failed at GetUserNameByID")
		return ""
	}
	return
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
	panic("Null user was entered")
}

/*
Hashes using the SHA256 package
*/
func HashString(text string) string {

	// Inits the hash
	h := sha256.New()

	// Bytes to string
	h.Write([]byte(text))

	// sum it
	z := h.Sum(nil)

	hashString := hex.EncodeToString(z)
	// Convert to string before returning
	return hashString
}

/*
* Called on process start. This will find the user if it exists in the database. If not, it will ask to create an account.
 */
func HandleAccountStartup() {

	// Gets the mac address
	macAddr, _ := getMacAddress()

	// Hashes that mac address to compare
	macAddr = HashString(macAddr)

	// searches for user using the hashed mac address
	row := conn.QueryRow("SELECT * FROM users WHERE macaddr = $1", macAddr)

	var id int
	var name, keyword string
	var macaddr string

	// Scans values into the variables
	err := row.Scan(&id, &name, &keyword, &macaddr)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Print("Error no row found for mac address!")
		}
	}

	// Checks for empty username
	if len(name) == 0 {
		// Prompt user for inputs
		fmt.Print("You appear to be a new user!\nEnter your username to get started: ")
		_, err := fmt.Scan(&name)
		if err != nil {
			os.Exit(1)
		}

		// Creates account using the hashed mac addresss and username that the user defined
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

	// Creates the friend code so that this user can be added as a friend
	code := generateFriendCode()

	// Inserts that data inside of the datbase
	_, err := conn.Exec("INSERT INTO users (name, keyword, macaddr) VALUES ($1, $2, $3)", username, code, macAddress)
	if err != nil {
		fmt.Printf("Failed at CreateAccount: %s", err)
	}
}

// Retrieves a user's freind code based on their name, which is passed in
func GetUserFriendCode(keyword string) (userKey string) {
	err := conn.QueryRow("SELECT keyword FROM users WHERE name=$1", keyword).Scan(&userKey)
	if err != nil {
		panic("Failed at GetUserFriendCode")
	}
	return
}

// Retrieves a user's id based on their name, which is passed in
func GetUserID(name string) (userID int) {
	err := conn.QueryRow("SELECT id FROM users WHERE name=$1", name).Scan(&userID)
	if err != nil {
		fmt.Print("Failed at GetUserID")
		return -1
	}
	return
}
