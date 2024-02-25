package main

import (
	"Example/src/database"
	"Example/src/network"
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
		pend:    *pend,
	}

	result := CheckInputs(flags)
	fmt.Println(result[0], result[1])
	if result[0] == "f" { // friend
		code := result[1]
		fmt.Println(code)
		friendsCode := database.GetUserFriendCode(code)
		result := database.AddFriend(friendsCode, curUserName)
		if result == false {
			fmt.Print("Failed to add friend! Not found!")
		} else {
			fmt.Print("Use has been added!")
		}

	} else if result[0] == "pend" {
		database.GetPendingTransfers()
	} else if result[0] == "r" { // retrieve

		// Receive message using result[1]
		transaction_identifier := "abcd 123"
		network.Receive_file(transaction_identifier)
	} else if result[0] == "d" { // delete friend

		// index, _ := strconv.Atoi(result[1])
		// Delete friend using result[1]

	} else if len(result) == 2 { // send
		// start daemon
		network.Send_file(result[0], result[1])
	} else {
		log.Fatal("No arguments given that match anything available")
	}

	fmt.Println(CheckInputs(flags))
}
