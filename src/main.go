package main

import (
	"Example/src/database"
	"flag"
	"fmt"
	"log"
	// "strconv"
)

// Main method for runnning the system
func main() {
	database.InitializeDatabase()
	database.HandleAccountStartup()
	user := database.GetCurrentUser()

	fmt.Printf("user: %v\n", user)
	// if user.name != "" {
	// 	fmt.Printf("Your code: %s\n", user.keyword)
	// }

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
	if result[0] == "f" { // friend
		fmt.Print("Friend command!")
		_, name, _, _ := database.GetUserDetails()
		database.AddFriend(result[1], name)
		// Add friend using result[1]
	} else if result[0] == "r" { // retrieve

		// index, _ := strconv.Atoi(result[1])

		// Receive message using result[1]
	} else if result[0] == "d" { // delete friend

		// index, _ := strconv.Atoi(result[1])
		// Delete friend using result[1]

	} else if len(result) == 2 { // send
		// this happens in a separate deamon thread
		// Send file
	} else {
		log.Fatal("No arguments given that match anything available")
	}

	fmt.Println(CheckInputs(flags))
}
