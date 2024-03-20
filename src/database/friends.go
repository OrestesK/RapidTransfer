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
func AddFriend(friendCode string, sender_id int) {
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

	isFriend := AreMutualFriends(sender_id, to_friend_id)
	if isFriend {
		log.Printf("User %s has already been added.\n", friend_name)
	} else {
		// New user is added
		_, err = conn.Exec(`
	INSERT INTO friends (user_one, user_two) VALUES ($1, $2)
	`, sender_id, to_friend_id)
		if err != nil {
			log.Printf("Error %s has occured.\n", err)
		}
	}

}

// DeleteFriend deletes a one-way friendship between two users.
func DeleteFriend(senderid int, recieverName string) {
	_, err := conn.Exec("DELETE FROM friends WHERE (user_one=$1 AND user_two=$2) OR user_one=$2 AND user_two=$1) ", senderid, GetUserID(recieverName))
	if err != nil {
		fmt.Println("Failed to remove friend", err)
	}
}

// IsFriend checks if there is a friendship between two users.
func IsFriend(user_one_id int, user_two_id int) bool {
	var temp int
	err := conn.QueryRow(`
	SELECT friends.id 
	FROM friends 
	INNER JOIN users ON friends.user_one=users.id AND friends.user_two=users.id 
	WHERE friends.user_one=%s AND friends.user_two=%s`, user_one_id, user_two_id).Scan(&temp)
	if err == sql.ErrNoRows || temp == 0 {
		return false
	}
	return true
}

// AreMutualFriends determines if two friends are mutual friends.
func AreMutualFriends(user_one_id int, user_two_id int) bool {
	return IsFriend(user_one_id, user_two_id) || IsFriend(user_two_id, user_one_id)
}

type Friend struct {
	Name       string
	FriendCode string
}

// GetFriendsList retrieves a list of friends for a given user.
func GetFriendsList(user_id int) (friendsList []Friend) {
	user_one := getUserOneFriends(user_id)
	user_two := getUserTwoFriends(user_id)

	friendsList = append(user_one, user_two...)
	return
}

func getUserOneFriends(userId int) (friendsList []Friend) {
	query := `
        SELECT users.name, users.friend_code
        FROM users
        JOIN friends ON users.id = friends.user_two
        WHERE friends.user_one = $1`

	rows, err := conn.Query(query, userId)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var friend Friend
		if err := rows.Scan(&friend.Name, &friend.FriendCode); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		friendsList = append(friendsList, friend)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}

	return friendsList
}

func getUserTwoFriends(userId int) (friendsList []Friend) {
	query := `
        SELECT ARRAY_AGG(name), ARRAY_AGG(friend_code)
        FROM users
        JOIN friends ON users.id = friends.user_one
        WHERE friends.user_two = $1
		GROUP BY friends.user_two`

	rows, err := conn.Query(query, userId)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var friend Friend
		if err := rows.Scan(&friend.Name, &friend.FriendCode); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		friendsList = append(friendsList, friend)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}

	return friendsList
}
