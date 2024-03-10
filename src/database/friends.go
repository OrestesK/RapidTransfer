package database

import (
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
func AddFriend(friendCode string, senderName string) (success bool) {
	var friendID int
	fmt.Printf("friendcode: %s\n", friendCode)
	err := conn.QueryRow("SELECT id FROM users WHERE keyword=$1;", friendCode).Scan(&friendID)
	fmt.Printf("FriendID From Code: %d\n", friendID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	senderID := GetUserID(senderName)

	row := conn.QueryRow("SELECT COUNT(*) FROM friends WHERE orig_user = $1 AND friend_id = $2", senderID, friendID)
	var count int
	row.Scan(&count)
	if count > 0 {
		fmt.Print("Friend is already added!")
		return true
	}
	_, err = conn.Exec(`
	INSERT INTO friends VALUES ($1, $2), ($3, $4)
	`, senderID, friendID, friendID, senderID)
	if err != nil {
		panic(err)
	}
	return true
}

// DeleteFriend deletes a one-way friendship between two users.
func DeleteFriend(senderName string, recieverName string) {
	_, err := conn.Exec("DELETE FROM friends WHERE orig_user=$1 AND friend_id=$2", GetUserID(senderName), GetUserID(recieverName))
	if err != nil {
		fmt.Println("Query failed", err)
	}
}

// IsFriend checks if there is a friendship between two users.
func IsFriend(userName1 string, userName2 string) (isFriend bool) {
	isFriend = false
	err1 := conn.QueryRow("SELECT id FROM friends WHERE user_from=$1, user_to=$2", GetUserID(userName1), GetUserID(userName2)).Scan(&isFriend)
	if err1 != nil {
		isFriend = true
		return
	}
	return
}

// AreMutualFriends determines if two friends are mutual friends.
func AreMutualFriends(userName1 string, userName2 string) (areMutuals bool) {
	areMutuals = false
	return IsFriend(userName1, userName2) && IsFriend(userName2, userName1)
}

// GetFriendsList retrieves a list of friends for a given user.
func GetFriendsList(username string) (friendsList [][]string) {
	userId := GetUserID(username)

	conn.QueryRow(`
	
	SELECT ARRAY_AGG(name), ARRAY_AGG(keyword)
	FROM users 
	
	INNER JOIN friends ON friends.friend_id=users.id
	WHERE orig_user=$1 
	
	GROUP BY friends.orig_user`, userId).Scan(&friendsList)
	return
}
