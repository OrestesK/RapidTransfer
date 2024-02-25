package main

import (
	"flag"
	"fmt"
)

// Creates the flags that are going to be used and assigns them values
func InitFlags() (*string, *string, *string, *string, *string, *string, *string, *string, *string) {
	s := flag.String("s", "", "Send to user")
	p := flag.String("p", "", "Path to file")
	f := flag.String("f", "", "Adding user to friends list")
	r := flag.String("r", "", "code of message receiving")
	d := flag.String("d", "", "Index of message deleting")
	df := flag.String("df", "", "Removes friend when given name")
	pn := flag.String("pn", "all", "Pending file transfers")
	fl := flag.String("fl", "all", "Retrieve friend list")
	c := flag.String("c", "self", "Retrieve personal friend code")

	return s, p, f, r, d, pn, fl, c, df
}

// Checks the flags for data
func CheckInputs(flags Flag) [2]string {
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
	if len(flags.recieve) != 0 {
		return [...]string{"r", flags.recieve}
	}
	// Checks to see if user is deleting a file from the inbox
	if len(flags.delete) != 0 {
		return [...]string{"d", flags.delete}
	}
	// Checks for usage of pending results command
	if len(flags.pend) != 0 {
		return [...]string{"pn", flags.pend}

	}
	// Checks to see if someone use the delete friend command
	if len(flags.dfriend) != 0 {
		return [...]string{"df", flags.dfriend}
	}
	// Checks for usage of the code command
	if len(flags.code) != 0 {
		return [...]string{"c", flags.code}

	}
	// Checks for the usage of friends list command
	if len(flags.fList) != 0 {
		return [...]string{"fl", flags.fList}
	}
	// If nothing is entered we exit the program
	fmt.Println("Exited")
	return result
}

// Creation of the flag struct
type Flag struct {
	send    string
	path    string
	friend  string
	recieve string
	delete  string
	dfriend string
	pend    string
	fList   string
	code    string
}
