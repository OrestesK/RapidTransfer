package main

import (
	"Example/src/database"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func execute_daemon(user_to string, filename string) {
	// execute daemon, runs in background independent of this
	cmd := exec.Command("go", "run", "src/daemon.go", user_to, filename)

	// dont worry about this
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Start()
	syscall.Setpgid(cmd.Process.Pid, cmd.Process.Pid)

	// saved pid to file
	tt := fmt.Sprintf("%d\n", cmd.Process.Pid)
	os.WriteFile("pid", []byte(tt), 0755)
}

// Main method for runnning the system
func main() {
	database.InitializeDatabase()
	database.HandleAccountStartup()
	curUser := database.GetCurrentUser()

	fmt.Println(curUser)
	_, curUserName, _, _ := database.GetUserDetails()
	fmt.Println(curUserName)

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
		friendsCode := database.GetUserFriendCode(result[1])
		result := database.AddFriend(friendsCode, curUserName)
		if result == false {
			fmt.Print("Failed to add friend! Not found!")
		} else {
			fmt.Print("Use has been added!")
		}

	} else if result[0] == "r" { // retrieve

		// index, _ := strconv.Atoi(result[1])

		// Receive message using result[1]
	} else if result[0] == "d" { // delete friend

		// index, _ := strconv.Atoi(result[1])
		// Delete friend using result[1]

	} else if len(result) == 2 { // send
		// Send file
		// HERE WE WILL START DAEMON
		execute_daemon(result[0], result[1])
	} else {
		log.Fatal("No arguments given that match anything available")
	}

	fmt.Println(CheckInputs(flags))
}
