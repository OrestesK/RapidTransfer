package database

import (
	"fmt"

	"github.com/jackc/pgx"
)

// GetPendingTransfers retrieves pending file transfer requests for a given username.
func GetPendingTransfers(username string) {
	// Retrieves user ID
	userID := GetUserID(username)

	// Retrieves pending requests from the database
	rows, _ := conn.Query("SELECT users.name as host, filename FROM transfer INNER JOIN users ON users.id = transfer.from_user WHERE to_user = $1", userID)

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
GetTransaction retrieves the names of two users who have had a transaction with each other.
It uses the specific keyword associated with the transaction to fetch user information.
*/
func GetTransaction(filename string) (names []string) {
	var from_user_id int
	var to_user_id int
	err := conn.QueryRow("SELECT from_user, to_user FROM transaction WHERE filename=$1", filename).Scan(&from_user_id, &to_user_id)
	if err != nil {
		panic("Failed at GetTransaction")
	}
	userFromName := GetUserNameByID(from_user_id)
	userToName := GetUserNameByID(to_user_id)
	names = []string{userFromName, userToName}
	return
}

// GetAddressFromTransactionPhrase retrieves the address associated with a transaction keyword.
func GetAddressFromTransactionPhrase(filename string) string {
	row := conn.QueryRow("SELECT ip_address FROM transfer WHERE filename = $1", filename)
	var address string
	err := row.Scan(&address)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("No results found")
			return ""
		}
	}
	return address
}

/*
UserCanViewTransaction checks if a user has the right to view a transaction.
It verifies the user's involvement in the transaction using their ID and the transaction keyword.
*/
func UserCanViewTransaction(userId int, filename string) bool {
	row := conn.QueryRow("SELECT COUNT(*) FROM transfer WHERE from_user = $1 OR to_user = $2 AND filename = $3", userId, userId, filename)
	var count int

	err := row.Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("No results found")
			return false
		}
	}
	return count > 0
}

// DeleteTransactionWithAddress deletes a transaction record based on the given address.
func DeleteTransactionWithAddress(address string) {
	_, err := conn.Exec("DELETE FROM transfer WHERE ip_address = $1", address)
	if err != nil {
		panic(err)
	}
}

/*
PerformTransaction allows two users to send files to each other.
It checks if the users are mutual friends and generates a unique phrase for the transaction.
*/
func PerformTransaction(user_from string, user_to string, address string, filename string) {
	fmt.Printf("from user: %s to user: %s\n", user_from, user_to)

	from_user_id := GetUserID(user_from)
	to_user_id := GetUserID(user_to)

	// Checks if users are mutual friends
	if AreMutualFriends(from_user_id, to_user_id) {
		// Inserts the transaction record into the database
		_, err := conn.Exec("INSERT INTO transfer (from_user, to_user, ip_address, filename) VALUES ($1,$2,$3,$4)", from_user_id, to_user_id, address, filename)
		if err != nil {
			fmt.Print("Failed at PerformTransaction")
			panic(err)
		}
	}
	fmt.Print("Failed to create transaction")
}
