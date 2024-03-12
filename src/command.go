package main

import (
	"Rapid/src/database"
	"Rapid/src/network"
	"fmt"
	"log"
	"strings"
)

// Help command to list information about what you can run
func help() {
	fmt.Print("-send user # Used to send file to user, must use -file path flag to specify the file\n")
	fmt.Print("-file path # Used to specify the path to the file you are sending, must be used with -send\n")
	fmt.Print("-add user_id # Used to add a friend, user_id is the id you retrieve when you use -info\n")
	fmt.Print("-inbox # Used to retrive information about files you have yet to accept\n")
	fmt.Print("-delete filename # Used to remove a file from your inbox\n")
	fmt.Print("-boot friend_id # Used to remove a friend from your friends list\n")
	fmt.Print("-recieve file # Used to accept a file being sent to you\n")
	fmt.Print("-friends # Used to list all of your friends and their friend id\n")
	fmt.Print("-info # Used to display your account information\n")
}

type Flag struct {
	flag  string
	input string
}

// Takes the command line input and returns a struct of flags and inputs
func retrieveFlags(command string) []Flag {
	var args []string
	var flags []Flag
	args = strings.Split(command, " ")

	for i, arg := range args {
		// If user did not enter any arguments, continue
		if strings.TrimSpace(arg) == "" {
			break
		}
		if arg[0] == '-' {
			// Found flag
			var temp Flag
			temp.flag = arg
			temp.input = ""
			// Checks for any flags that actually need input. If they do, then add the next arg and increment i to skip it
			if strings.Compare(arg, "-send") == 0 || strings.Compare(arg, "-file") == 0 || strings.Compare(arg, "-add") == 0 || strings.Compare(arg, "-delete") == 0 || strings.Compare(arg, "-boot") == 0 || strings.Compare(temp.flag, "-recieve") == 0 {
				if len(args[i+1]) == 0 {
					log.Fatalf("Flag %s takes a value\n", temp.flag)
				}
				if args[i+1][0] == '-' {
					log.Fatalf("Flag %s takes a value\n", temp.flag)
				}
				temp.input = args[i+1]
				i++
			}
			//fmt.Printf("Value: %s Input: %s\n", arg, temp.input)

			flags = append(flags, temp)
		}

	}
	return flags

}

// Goes and runs through all the commands that are entered
func command(flags []Flag, user string) {
	sent := false
LOOP:
	for _, temp := range flags {
		switch argument := temp.flag; argument {
		case "-send":
			if sent {
				continue
			}
			for _, search := range flags {
				if strings.Compare(search.flag, "-file") == 0 {
					network.Send_file(temp.input, search.input)
					fmt.Println("File has been sent and will be waiting to be accepted")
					sent = true
				}
			}
			if !sent {
				log.Fatal("Need to specify a file to send\n")
			}
			break LOOP
		case "-file":
			if sent {
				continue
			}
			for _, search := range flags {
				if strings.Compare(search.flag, "-send") == 0 {
					network.Send_file(temp.input, search.input)
					fmt.Println("File has been sent and is waiting to be accepted")
					sent = true
				}
			}
			if !sent {
				log.Fatal("Need to specify a person to send the file to\n")
			}
			break LOOP
		case "-add":
			database.AddFriend(temp.input, user)
			fmt.Println("Inbox has been displayed")
			break LOOP
		case "-inbox":
			database.GetPendingTransfers(user)
			fmt.Println("Inbox has been displayed")
			break LOOP
		case "-delete":
			network.DeleteFile(temp.input)
			fmt.Println("File has been deleted")
			break LOOP
		case "-boot":
			database.DeleteFriend(user, temp.input)
			fmt.Println("Friend has been deleted")
			break LOOP
		case "-recieve":
			network.RecieveFile(temp.input)
			fmt.Println("File has been received")
			break LOOP
		case "-friends":
			friendList := database.GetFriendsList(user)
			for _, temp := range friendList {
				fmt.Println("Friend name: ", temp)
			}
			fmt.Println("Friends list has been displayed")
			break LOOP
		case "-info":
			fmt.Printf("Username: %s, Friend_ID: %s\n", user, database.GetUserFriendCode(user))
			break LOOP
		case "-help":
			help()
			break LOOP
		default:
			log.Fatalf("Flag %s does not exist\n", argument)
			break LOOP
		}
	}
}
