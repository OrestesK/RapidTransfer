package database

import (
	"database/sql"
	"fmt"
	"log"
)

// GetUserNameByFriendCode retrieves a user's name based on their friend code.
func GetUserNameByFriendCode(friendCode string) (userName string) {
	err := conn.QueryRow("SELECT name FROM user WHERE id=$1", friendCode).Scan(&userName)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// AddFriend adds a friend to a sender's friends list based on their friend code.
func AddFriend(friendCode string, senderName string) {
	var to_friend_id int
	var friend_name string
	// Finds the friend code
	err := conn.QueryRow("SELECT id, name FROM users WHERE friend_code=$1;", friendCode).Scan(&to_friend_id, &friend_name)

	// Throws error if user does not exist
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User with the friend code %s does not exist.\n", friendCode)
		}
	}

	from_user_id := GetUserID(senderName)

	isFriend := AreMutualFriends(from_user_id, to_friend_id)
	if isFriend {
		log.Printf("User %s has already been added.\n", friend_name)
	} else {
		// New user is added
		_, err = conn.Exec(`
	INSERT INTO friends (user_one, user_two) VALUES ($1, $2)
	`, from_user_id, to_friend_id)
		if err != nil {
			log.Printf("Error %s has occured.\n", err)
		}
	}

}

// DeleteFriend deletes a one-way friendship between two users.
func DeleteFriend(senderName string, recieverName string) {
	_, err := conn.Exec("DELETE FROM friends WHERE user_one=$1 AND user_two=$2", GetUserID(senderName), GetUserID(recieverName))
	if err != nil {
		fmt.Println("Query failed", err)
	}
}

// IsFriend checks if there is a friendship between two users.
func IsFriend(user_one_id int, user_two_id int) bool {
	row := conn.QueryRow(`
	SELECT id 
	FROM friends 
	INNER JOIN users ON friends.user_one=users.id AND friends.user_two=users.id 
	WHERE user_from=$1 
	AND user_to=$2`, user_one_id, user_two_id)
	if row.Scan() == sql.ErrNoRows {
		return false
	}
	return true
}

// AreMutualFriends determines if two friends are mutual friends.
func AreMutualFriends(user_one_id int, user_two_id int) bool {
	return IsFriend(user_one_id, user_two_id) || IsFriend(user_two_id, user_one_id)
}

// GetFriendsList retrieves a list of friends for a given user.
func GetFriendsList(username string) (friendsList []string) {
	userId := GetUserID(username)

	conn.QueryRow(`
	
	SELECT name, friend_code
		FROM users
		JOIN friends ON users.id = friends.user_one OR users.id = friends.user_two
		WHERE friends.user_one = $1 OR friends.user_two = $1

		GROUP BY users.id`, userId).Scan(&friendsList)
	return
}
