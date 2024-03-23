package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
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
	z := h.Sum(nil)

	hashString := hex.EncodeToString(z)
	// Convert to string before returning
	return hashString
}

/*
Creates an account in the database with specifid username/uuid.
uuid is unique to the computer, and is used on startup to indentify the pc.
*/
func CreateAccount(username string, password string) error {

	// Hashes password
	password = HashInfo(password)

	// Retrievs UUID and hashes it
	uuid, err := getUUID()
	if err != nil {
		return err
	}
	uuid = HashInfo(uuid)
	if alreadyExistsCheck(username, uuid) == 0 {
		return errors.New("User already exist")
	}
	// Creates the friend code so that this user can be added as a friend
	code := generateFriendCode()

	// Inserts that data inside of the datbase
	_, err = conn.Exec("INSERT INTO users (name, password, friend_code, uuid) VALUES ($1, $2, $3, $4)", username, password, code, uuid)
	if err != nil {
		return err
	}
	return nil
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
func SetCurrentUsersId(name string, password string) error {
	// Retrievs UUID and hashes it
	uuid, err := getUUID()
	if err != nil {
		return err
	}
	uuid = HashInfo(uuid)
	password = HashInfo(password)
	row := conn.QueryRow(
		`SELECT id 
		FROM users 
		WHERE name=$1 AND password=$2 AND uuid=$3`, name, password, uuid)

	err = row.Scan(&current_user)

	if err != nil {
		return errors.New("Username or password is wrong\n")
	}
	return nil
}

func GetCurrentId() int {
	return current_user
}
