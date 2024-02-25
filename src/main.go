package main

import (
	"Example/src/database"
	"flag"
	"fmt"
	"log"
)

// Main method for runnning the system
func main() {
	database.InitializeDatabase()
	database.HandleAccountStartup()
	curUser := database.GetCurrentUser()

	fmt.Println(curUser)
	_, curUserName, _, _ := database.GetUserDetails()
	fmt.Println(curUserName)

	s, p, friend, r, d, pend := InitFlags()
	flag.Parse()
	flags := Flag{
		send:    *s,
		path:    *p,
		friend:  *friend,
		recieve: *r,
		delete:  *d,
		pend: *pend,
	}

	result := CheckInputs(flags)
	if result[0] == "f" { // friend
		friendsCode := database.GetUserFriendCode(result[1])
		result := database.AddFriend(friendsCode, curUserName)
		if (result == false) {
			fmt.Print("Failed to add friend! Not found!")
		} else {
			fmt.Print("Use has been added!")
		}

	} else if result[0] == "pend" {
		database.GetPendingTransfers()
	}else if result[0] == "r" { // retrieve

		// index, _ := strconv.Atoi(result[1])

		// Receive message using result[1]
	} else if result[0] == "d" { // delete friend

		// index, _ := strconv.Atoi(result[1])
		// Delete friend using result[1]

	} else if len(result) == 2 { // send
		// database.PerformTransaction()
		// Send file
	} else {
		log.Fatal("No arguments given that match anything available")
	}

	fmt.Println(CheckInputs(flags))
}
