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

	_, curUserName, _, _ := database.GetUserDetails()

	s, p, f, r, d, pn, fl := InitFlags()
	flag.Parse()
	flags := Flag{
		send:    *s,
		path:    *p,
		friend:  *f,
		recieve: *r,
		delete:  *d,
		fList:   *fl,
		pend:    *pn,
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

	} else if result[0] == "pn" {
		database.GetPendingTransfers()
		friends := database.GetFriendsList(curUserName)
		fmt.Printf("friends: %v\n", friends)
	} else if result[0] == "r" { // retrieve

		// Receive message using result[1]
		network.Receive_file(result[1])
	} else if result[0] == "d" { // delete friend

		// index, _ := strconv.Atoi(result[1])
		// Delete friend using result[1]

	} else if result[0] == "fl" {
		friendList := database.GetFriendsList(curUserName)
		for namez := range friendList {
			fmt.Println("Friend name: ", namez)
		}
	} else if len(result) == 2 { // send
		// start daemon
		network.Send_file(result[0], result[1])
	} else {
		log.Fatal("No arguments given that match anything available")
	}

	// fmt.Println(CheckInputs(flags))
}
