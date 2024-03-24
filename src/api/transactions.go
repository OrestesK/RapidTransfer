package database

import custom "github.com/Zaikoa/rapid/src/handling"

type Transaction struct {
	From_user string
	File_name string
}

/*
retrieves pending file transfer requests for a certain user
*/
func GetPendingTransfers(id int) ([]Transaction, error) {
	if id == 0 {
		return nil, custom.NewError("User must be logged in to use this method")
	}

	query := `
	SELECT nickname, filename
	FROM transfer 
	INNER JOIN users ON transfer.from_user = users.id
	WHERE to_user=$1`
	rows, err := conn.Query(query, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Formats all the transaction details and appends to struct list
	var inbox []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.From_user, &transaction.File_name); err != nil {
			return nil, err
		}
		inbox = append(inbox, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inbox, nil
}

/*
Checks if a user has the right to view a transaction
*/
func UserCanViewTransaction(id int, filename string) bool {
	query := `SELECT id FROM transfer WHERE (from_user = $1 OR to_user = $2) AND filename = $3`
	row := conn.QueryRow(query, id, id, filename)

	// dummy value that will store the id and is not checked (just used to check if no rows)
	var temp int // Temp cannot be null since sql does not support null primary keys
	err := row.Scan(&temp)
	if err != nil || temp == 0 {
		return false
	}

	return true
}

/*
Deletes a transaction record based on the given id
*/
func DeleteTransaction(key string) error {
	query := `DELETE FROM transfer WHERE key = $1`
	_, err := conn.Exec(query, key)
	if err != nil {
		return err
	}
	return nil
}

/*
Enters the information of a transaction to the database for later referal
*/
func PerformTransaction(from_user_id int, user_to string, filename string, key string) (bool, error) {
	to_user_id, _ := GetUserID(user_to)

	if AreMutualFriends(from_user_id, to_user_id) {
		query := `INSERT INTO transfer (from_user, to_user, key, filename) VALUES ($1,$2,$3,$4)`
		_, err := conn.Exec(query, from_user_id, to_user_id, key, filename)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

/*
Retrieves decription key from the database that is used to decrypt the file
*/
func RetrieveKey(filename string, to_user int) (string, error) {
	var key string
	query := `SELECT key FROM transfer WHERE to_user = $1 AND filename = $2`
	err := conn.QueryRow(query, to_user, filename).Scan(&key)

	if err != nil {
		return "", err
	}

	return key, nil
}
