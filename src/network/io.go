package network

import (
	"Rapid/src/database"
	"fmt"

	"os/exec"
)

func Send_file(user_from string, user_to string, filename string) {
	// execute daemon, runs in background independent of this
	cmd := exec.Command("go", "run", "src/daemon/daemon.go", user_from, user_to, filename)

	cmd.Start()

	// saved pid to file
	// tt := fmt.Sprintf("%d\n", cmd.Process.Pid)
	// os.WriteFile("pid", []byte(tt), 0755)
}

func RecieveFile(user string, filename string) {
	node := Initialize_node()
	id := database.GetUserID(user)
	result := database.UserCanViewTransaction(id, filename)

	// get big key from small key
	address, err := database.GetAddressFromFileName(filename)
	done := make(chan bool)
	// client
	if (err == nil && address == "") || !result {
		fmt.Println("File Not found")
	} else {
		Client(node, address, filename, done, false)
		<-done
	}
}

func DeleteFile(filename string) {
	node := Initialize_node()
	_, name, _, _ := database.GetUserDetails()
	id := database.GetUserID(name)
	result := database.UserCanViewTransaction(id, filename)

	// get big key from small key
	address, err := database.GetAddressFromFileName(filename)
	done := make(chan bool)
	// client
	if (err == nil && address == "") || !result {
		fmt.Println("File Not found")
	} else {
		Client(node, address, filename, done, true)
		<-done
	}
}
