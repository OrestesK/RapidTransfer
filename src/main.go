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

	s, p, friend, r, d := InitFlags()
	flag.Parse()
	flags := Flag{
		send:    *s,
		path:    *p,
		friend:  *friend,
		recieve: *r,
		delete:  *d,
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

	} else if result[0] == "r" { // retrieve

	} else if result[0] == "d" { // delete file

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
