package database

import (
	custom "Rapid/src/handling"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"math/rand"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Stores the current users id
var current_user int

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
			return "", err
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
			return "", err
		} else {
			return out, nil
		}
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("findmnt", "/", "-o", "UUID", "-n")
		var b []byte
		b, err := cmd.CombinedOutput()
		out := string(b)

		if err != nil {
			return "", err
		} else {
			return out, nil
		}
	} else {
		return "", nil
	}
}

/*
Retrieves a user's name based on their id, which is passed in
*/
func GetUserNameByID(id int) (name string) {
	query := `SELECT name FROM users WHERE id=$1`
	conn.QueryRow(query, id).Scan(&name)
	return
}

/*
Generates a random friend code that will be assigned to a user during account creation
*/
func generateFriendCode() string {
	rand.Seed(time.Now().UnixNano())
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
	h := sha256.New()
	h.Write([]byte(text))
	z := h.Sum(nil)
	hashString := hex.EncodeToString(z)

	return hashString
}

/*
Creates an account in the database
uuid is unique to the computer, and is used on startup to indentify the device
*/
func CreateAccount(username string, password string) error {
	password = HashInfo(password)
	uuid, err := getUUID()
	if err != nil {
		return err
	}
	uuid = HashInfo(uuid)
	result, _ := alreadyExistsCheck(username, uuid)
	if result == 0 {
		return custom.NewError("User already exist")
	}
	code := generateFriendCode()

	// Inserts that data inside of the datbase
	query := `INSERT INTO users (name, password, friend_code, uuid) VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(query, username, password, code, uuid)
	if err != nil {
		return err
	}
	return nil
}

/*
Retrieves a user's freind code based on their name, which is passed in
*/
func GetUserFriendCode(id int) (string, error) {
	var code string
	query := `SELECT friend_code FROM users WHERE id=$1`
	err := conn.QueryRow(query, id).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}

/*
Retrieves a user's id based on their name, which is passed in
*/
func GetUserID(name string) (int, error) {
	var id int
	query := `SELECT id FROM users WHERE name=$1`
	err := conn.QueryRow(query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return id, custom.NewError("User not found")
		}
		return id, err
	}
	return id, nil
}

/*
If user does not exist, the id returns 0, this is because sql cannot have a null primary id
*/
func alreadyExistsCheck(name string, uuid string) (int, error) {
	var id int
	query := `SELECT id FROM users WHERE name=$1 AND uuid=$2`
	err := conn.QueryRow(query, name, uuid).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

/*
Sets the user who is currently logged in
*/
func SetCurrentUsersId(name string, password string) error {
	uuid, err := getUUID()
	if err != nil {
		return err
	}
	uuid = HashInfo(uuid)
	password = HashInfo(password)
	query := `SELECT id FROM users WHERE name=$1 AND password=$2 AND uuid=$3`
	err = conn.QueryRow(query, name, password, uuid).Scan(&current_user)
	if err != nil {
		return custom.NewError("Username or password is wrong")
	}

	return nil
}

// Returns the current users id
func GetCurrentId() int {
	return current_user
}
