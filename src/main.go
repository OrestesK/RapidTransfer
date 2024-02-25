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

		// Add friend using result[1]
	} else if result[0] == "r" { // retrieve

		// index, _ := strconv.Atoi(result[1])

		// Receive message using result[1]
	} else if result[0] == "d" { // delete friend

		// index, _ := strconv.Atoi(result[1])
		// Delete friend using result[1]

	} else if len(result) == 2 { // send
		// Send file
	} else {
		log.Fatal("No arguments given that match anything available")
	}

	fmt.Println(CheckInputs(flags))
}
