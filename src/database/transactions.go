package database

import (
	"fmt"
	"log"

	"github.com/jackc/pgx"
)

// GetPendingTransfers retrieves pending file transfer requests for a given user id.
func GetPendingTransfers(user_id int) {

	// Retrieves pending requests from the database
	rows, _ := conn.Query("SELECT from_user, filename FROM transfer WHERE to_user=$1", user_id)

	// Prints out information for each pending transfer
	for rows.Next() {
		var host, filename string
		err := rows.Scan(&host, &filename)

		// Handles the error while scanning the row
		if err != nil {
			fmt.Println("Error scanning row:", err)
			break
		}
		fmt.Printf("From: %s Filename: %s\n", host, filename)
	}
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
func DeleteTransaction(id string) {
	// Deletes from table
	_, err := conn.Exec("DELETE * FROM transfer WHERE id = $1", id)

	// Logs the error
	if err != nil {
		log.Fatal(err)
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
			log.Fatalf("Error inserting transfer data: %s", err)
		}
		return true
	}
	return false
}

func RetrieveKey(filename string, to_user int) string {
	row := conn.QueryRow("SELECT key FROM transfer WHERE to_user = $2 AND filename = $3", to_user, filename)

	// dummy value that will store the id and is not checked (just used to check if no rows)
	var key string
	err := row.Scan(&key)
	if err != nil {
		log.Fatal(err)
	}

	return key
}
