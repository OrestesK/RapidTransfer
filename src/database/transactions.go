package database

import (
	"fmt"

	"github.com/jackc/pgx"
)

type Transaction struct {
	From_user string
	File_name string
}

// GetPendingTransfers retrieves pending file transfer requests for a given user id.
func GetPendingTransfers(user_id int) (inbox []Transaction) {

	// Retrieves pending requests from the database
	rows, err := conn.Query("SELECT from_user, filename FROM transfer WHERE to_user=$1", user_id)

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.From_user, &transaction.File_name); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		inbox = append(inbox, transaction)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}

	return inbox
}

/*
UserCanViewTransaction checks if a user has the right to view a transaction.
It verifies the user's involvement in the transaction using their ID and the key of a transaction.
*/
func UserCanViewTransaction(userId int, filename string) bool {
	row := conn.QueryRow("SELECT id FROM transfer WHERE (from_user = $1 OR to_user = $2) AND filename = $3", userId, userId, filename)

	// dummy value that will store the id and is not checked (just used to check if no rows)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Then user cannot view transaction
			return false
		}
	}
	// User can see transaction
	return true
}

// Deletes a transaction record based on the given id.
func DeleteTransaction(key string) error {
	// Deletes from table
	_, err := conn.Exec("DELETE FROM transfer WHERE key = $1", key)

	// Logs the error
	if err != nil {
		return err
	}
}

/*
PerformTransaction allows two users to send files to each other.
*/
func PerformTransaction(from_user_id int, user_to string, filename string, key string) bool {
	to_user_id := GetUserID(user_to)

	// Checks if users are mutual friends
	if AreMutualFriends(from_user_id, to_user_id) {
		// Inserts the transaction record into the database
		_, err := conn.Exec("INSERT INTO transfer (from_user, to_user, key, filename) VALUES ($1,$2,$3,$4)", from_user_id, to_user_id, key, filename)
		if err != nil {
			fmt.Println(err)
			return false
		}
		return true
	}
	return false
}

func RetrieveKey(filename string, to_user int) (string, error) {
	row := conn.QueryRow("SELECT key FROM transfer WHERE to_user = $1 AND filename = $2", to_user, filename)

	// dummy value that will store the id and is not checked (just used to check if no rows)
	var key string
	err := row.Scan(&key)
	if err != nil {
		return "", err
	}

	return key, nil
}
