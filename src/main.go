package main

import (
	"Rapid/src/database"
	"fmt"
	"os"
	"strings"
)

// Main method for runnning the system
func main() {
	database.InitializeDatabase()
	database.HandleAccountStartup()

	_, user, _, _ := database.GetUserDetails()
	fmt.Printf("Currently Logged in as %s\n", user)

	// Creates a splice from the command line input
	splice := os.Args[1:]
	arguments := strings.Join(splice, " ")

	// Retreives flags called and then runs the commands inside of them
	flags := retrieveFlags(arguments)
	command(flags, user)
}
