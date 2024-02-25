package main

import (
	"Example/src/database"
	"Example/src/network"
	"flag"
	"fmt"
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
	switch argument := result[0]; argument {
	case "f":
		code := result[1]
		fmt.Println(code)
		friendsCode := database.GetUserFriendCode(code)
		result := database.AddFriend(friendsCode, curUserName)
		if result == false {
			fmt.Print("Failed to add friend! Not found!")
		} else {
			fmt.Print("Use has been added!")
		}

	case "pn":
		database.GetPendingTransfers()
		friends := database.GetFriendsList(curUserName)
		fmt.Printf("friends: %v\n", friends)

	case "r":
		network.Receive_file(result[1])

	case "d":

	case "fl":
		friendList := database.GetFriendsList(curUserName)
		for namez := range friendList {
			fmt.Println("Friend name: ", namez)
		}
	default:
		network.Send_file(result[0], result[1])
	}

}
