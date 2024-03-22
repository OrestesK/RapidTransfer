package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

/*
Returns users UUID by first checking what system they are on
*/
func getUUID() (string, error) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "path", "win32_computersystemproduct", "get", "UUID")
		var b []byte
		b, err := cmd.CombinedOutput()
		out := string(b)

		if err != nil {
			return "", nil
		} else {
			result := strings.Split(out, "\n")
			return result[1], nil
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("uuidgen")
		var b []byte
		b, err := cmd.CombinedOutput()
		out := string(b)

		if err != nil {
			return "", nil
		} else {
			return out, nil
		}
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("findmnt", "/", "-o", "UUID", "-n")
		var b []byte
		b, err := cmd.CombinedOutput()
		out := string(b)

		if err != nil {
			return "", nil
		} else {
			return out, nil
		}
	} else {
		return "", nil
	}
}

var current_user int

// Retrieves a user's name based on their id, which is passed in
func GetUserNameByID(id int) (userName string) {
	conn.QueryRow("SELECT name FROM users WHERE id=$1", id).Scan(&userName)
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

/*
Hashes using the SHA256 package
*/
func HashInfo(text string) string {

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
	// Inits all the information
	var name, password string
	// Gets the uuid address
	uuid, err := getUUID()
	if err != nil {
		fmt.Println(err)
	}

	// Hashes that uuid to compare
	uuid = HashInfo(uuid)

	fmt.Print("Enter your username and password to login or create to create a new user\nExpected: username password or create\n")

	// May look weird but it does check if the input was create and if it isnt, just stores the username
	_, err = fmt.Scan(&name)
	if strings.Compare(name, "create") == 0 {
		// User entered create
		fmt.Println("You have decided to create a user, if you have other users they will still be accesible")
		fmt.Println("Enter your username and password to get started\nExpected: username password")
		_, err := fmt.Scan(&name, &password)

		// Checks to make sure we cannot create multiple users on the same device with the same name
		for alreadyExistsCheck(name, uuid) != -1 {
			fmt.Println("Duplicate creation of user on one machine caught. Change username to procceed")
			_, err = fmt.Scan(&name)
		}
		if err != nil {
			os.Exit(1)
		}
		password = HashInfo(password)
		createAccount(name, password, uuid)
	} else {
		_, err = fmt.Scan(&password)
		//fmt.Println(password)
		password = HashInfo(password)
		if err != nil {
			os.Exit(1)
		}
	}
	err = setCurrentUsersId(name, password, uuid)
	if err != nil {
		log.Fatal(err)
	}
}

/*
Creates an account in the database with specifid username/uuid.
uuid is unique to the computer, and is used on startup to indentify the pc.
*/
func createAccount(username string, password string, uuid string) {

	// Creates the friend code so that this user can be added as a friend
	code := generateFriendCode()

	// Inserts that data inside of the datbase
	_, err := conn.Exec("INSERT INTO users (name, password, friend_code, uuid) VALUES ($1, $2, $3, $4)", username, password, code, uuid)
	if err != nil {
		fmt.Printf("Failed to create account", err)
	}
}

// Retrieves a user's freind code based on their name, which is passed in
func GetUserFriendCode(user_id int) (friend_code string) {
	err := conn.QueryRow("SELECT friend_code FROM users WHERE id=$1", user_id).Scan(&friend_code)
	if err != nil {
		panic("Failed at GetUserFriendCode")
	}
	return
}

// Retrieves a user's id based on their name, which is passed in
func GetUserID(name string) (userID int) {
	err := conn.QueryRow("SELECT id FROM users WHERE name=$1", name).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			return -1
		}
		fmt.Println("Error retrieving user ID:", err)
		return -1
	}
	return userID
}

// Retrieves a user's id based on their name, which is passed in
func alreadyExistsCheck(name string, uuid string) (userID int) {
	err := conn.QueryRow("SELECT id FROM users WHERE name=$1 AND uuid=$2", name, uuid).Scan(&userID)
	if err != nil {
		return -1
	}
	return
}
func setCurrentUsersId(name string, password string, uuid string) error {
	row := conn.QueryRow(
		`SELECT id 
		FROM users 
		WHERE name=$1 AND password=$2 AND uuid=$3`, name, password, uuid)

	err := row.Scan(&current_user)

	if err != nil {
		return errors.New("Username or password is wrong\n")
	}
	return nil
}

func GetCurrentId() int {
	return current_user
}
