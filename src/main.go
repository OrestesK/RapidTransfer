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

	s, p, f, r, d, pn, fl, c := InitFlags()
	flag.Parse()
	flags := Flag{
		send:    *s,
		path:    *p,
		friend:  *f,
		recieve: *r,
		delete:  *d,
		fList:   *fl,
		pend:    *pn,
		code:    *c,
	}

	result := CheckInputs(flags)
	switch argument := result[0]; argument {
	case "f":
		code := result[1]
		fmt.Println(code)
		friendsCode := database.GetUserFriendCode(code)
		result := database.AddFriend(friendsCode, curUserName)
		if !result {
			fmt.Println("Failed to add friend! Not found!")
		} else {
			fmt.Println("Use has been added!")
		}

	case "pn":
		database.GetPendingTransfers(curUserName)
	case "r":
		network.Receive_file(result[1])

	case "d":

	case "fl":
		friendList := database.GetFriendsList(curUserName)
		for _, namez := range friendList {
			fmt.Println("Friend name: ", namez)
		}
	case "c":
		fmt.Println(database.GetUserFriendCode(curUserName))
	default:
		network.Send_file(result[0], result[1])
	}

}
