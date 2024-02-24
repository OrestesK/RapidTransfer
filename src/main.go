package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
)

// Creates the flags that are going to be used and assigns them values
func initFlags() (*string, *string, *string, *int, *int) {
	s := flag.String("s", "", "Send to user")
	p := flag.String("p", "", "Path to file")
	friend := flag.String("friend", "", "Adding user to friends list")
	r := flag.Int("r", -1, "Index of message receiving")
	d := flag.Int("d", -1, "Index of message deleting")
	return s, p, friend, r, d
}

// Checks the flags for data
func checkInputs(flags Flag) [2]string {
	var result [2]string
	// Checks to see if the send flag was used
	if len(flags.send) != 0 && len(flags.path) != 0 {
		// Formats the send and file path arguments
		result := [...]string{flags.send, flags.path}
		return result
	}
	// Checks if the user is adding a friend
	if len(flags.friend) != 0 {
		return [...]string{"f", flags.friend}
	}
	// Checks to see if user is receiving a file from inbox
	if flags.recieve != -1 {
		return [...]string{"r", string(flags.recieve)}
	}
	// Checks to see if user is deleting a file from the inbox
	if flags.delete != -1 {
		return [...]string{"d", string(flags.delete)}
	}
	// If nothing is entered we exit the program
	fmt.Println("Exited")
	return result
}

type Flag struct {
	send    string
	path    string
	friend  string
	recieve int
	delete  int
}

// Main method for runnning the system
func main() {
	s, p, friend, r, d := initFlags()
	flag.Parse()
	flags := Flag{
		send:    *s,
		path:    *p,
		friend:  *friend,
		recieve: *r,
		delete:  *d,
	}

	result := checkInputs(flags)
	if result[0] == "f" {
		// Add friend using result[1]
	} else if result[0] == "r" {
		index, _ := strconv.Atoi(result[1])
		// Receive message using result[1]
	} else if result[0] == "d" {
		index, _ := strconv.Atoi(result[1])
		// Delete friend using result[1]
	} else if len(result) == 2 {
		// Send file
	} else {
		log.Fatal("No arguments given that match anything available")
	}

	fmt.Println(checkInputs(flags))