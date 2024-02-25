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
	// Retrieves the flags from the init
	s, p, f, r, d, pn, fl, c, df := InitFlags()
	flag.Parse()

	// Creation of the the flag struct and all flags that can be called
	flags := Flag{
		send:    *s,
		path:    *p,
		friend:  *f,
		recieve: *r,
		delete:  *d,
		dfriend: *df,
		fList:   *fl,
		pend:    *pn,
		code:    *c,
	}

	// Checks the flags and sees which ones are used and valid for calling
	result := CheckInputs(flags)
	fmt.Println(result[0], result[1])
	switch argument := result[0]; argument {
	// Adds friend to your friends list, Usage -f name
	case "f": //works
		code := result[1]
		fmt.Println(code)
		result := database.AddFriend(code, curUserName)
		if !result {
			fmt.Println("Failed to add friend! Not found!")
		} else {
			fmt.Println("User has been added!")
		}
	// Retrieves all of the pending transfers, Usage -pn all
	case "pn": // Works
		fmt.Println("This is a debugging statement")
		database.GetPendingTransfers(curUserName)
	case "r": // Works
		network.Receive_file(result[1])
		fmt.Println("File has been received")
	// Deletes file inside of the inbox, usage -d index
	case "d": // Works
		network.Fake_receive_file(result[1])
		fmt.Println("File has been deleted")
	// Deletes friend when given the username
	case "df": // Does not work
		database.DeleteFriend(curUserName, result[1])
		fmt.Println("Friend has been deleted")
	// Retrieves the users friend list, usage -fl all
	case "fl": // works
		friendList := database.GetFriendsList(curUserName)
		//fmt.Println("This is a debugging statement")
		for _, namez := range friendList {
			fmt.Println("Friend name: ", namez)
		}
	// Retrieves the users friend code, Usage -c self
	case "c": // works
		fmt.Println(database.GetUserFriendCode(curUserName))
	default: // works
		// Sending file to user, usage to_user file_path
		network.Send_file(result[0], result[1])
		fmt.Println("File has been sent and will be waiting to be accepted")
	}

}
