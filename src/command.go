package main

import (
	"Rapid/src/cloud"
	"Rapid/src/database"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
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
				if len(args[i+1]) == 0 || args[i+1][0] == '-' {
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
func command(flags []Flag, user int) {
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
					err, result := cloud.UploadToMega(search.input, user, temp.input)
					if err != nil {
						fmt.Println(err)
					}
					if result {
						fmt.Println("File has been sent and will be waiting to be accepted")
					} else {
						fmt.Println("The requested user either does not exist or is not added")
					}
					sent = true
				}
			}
			if !sent {
				fmt.Println("Need to specify a file to send")
			}
			break LOOP
		case "-file":
			if sent {
				continue
			}
			for _, search := range flags {
				if strings.Compare(search.flag, "-send") == 0 {
					err, result := cloud.UploadToMega(search.input, user, temp.input)
					if err != nil {
						fmt.Println(err)
					}
					if result {
						fmt.Println("File has been sent and will be waiting to be accepted")
					} else {
						fmt.Println("The requested user either does not exist or is not added")
					}
					sent = true
				}
			}
			if !sent {
				fmt.Println("Need to specify a person to send the file to")
			}
			break LOOP
		case "-add":
			result, err := database.AddFriend(temp.input, user)
			if err != nil {
				fmt.Println(err)
			}

			if !result {
				fmt.Println("There does not exist a user with that friend code.")

			}

			if result {
				fmt.Println("Friend has been added")
			}

			break LOOP
		case "-inbox":
			displayInbox(user)
			fmt.Println("Inbox has been displayed")
			break LOOP
		case "-delete":
			result, err := cloud.DeleteFromMega(user, temp.input)
			if err != nil {
				fmt.Println(err)
			}
			if result {
				fmt.Println("File has been deleted")
			} else {
				fmt.Println("Could not delete the file from the inbox")
			}
			break LOOP
		case "-boot":
			result, err := database.DeleteFriend(user, temp.input)
			if err != nil {
				fmt.Println("Failed to remove friend", err)
			}
			if result {
				fmt.Println("Friend has been deleted")
			}
			break LOOP
		case "-recieve":
			_, result := cloud.DownloadFromMega(user, temp.input, "")
			if result {
				fmt.Println("File has been received")
			} else {
				fmt.Println("Filename or item does not not exist within your inbox")
			}
			break LOOP
		case "-friends":
			displayFriends(user)
			fmt.Println("Friends list has been displayed")
			break LOOP
		case "-info":
			fmt.Printf("| Username   %s | Friend code   %s |\n", database.GetUserNameByID(user), database.GetUserFriendCode(user))
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

// displays friends list
func displayFriends(user int) {
	friendsList := database.GetFriendsList(user)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Friend Code"})
	for _, friend := range friendsList {
		t.AppendRows([]table.Row{
			{friend.Name, friend.FriendCode},
		})
	}
	t.Render()
}

// displays inbox
func displayInbox(user int) {
	inbox := database.GetPendingTransfers(user)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"From", "File Name"})
	for _, transaction := range inbox {
		t.AppendRows([]table.Row{
			{transaction.From_user, transaction.File_name},
		})
	}
	t.Render()
}
