package main

import (
	"Example/src/database"
	"Example/src/network"
	"os"
)

// this is the daemon
func main() {
	args := os.Args[1:]
	user_to := args[0]
	file_name := args[1]

	os.WriteFile("to", []byte(user_to), 0755)
	os.WriteFile("file_name", []byte(file_name), 0755)

	node := network.Initialize_node()

	done := make(chan bool)
	// I will computer and provide key
	address := network.Server(node, done)

	// initialize user
	database.HandleAccountStartup()
	_, user_from, _, _ := database.GetUserDetails()

	database.PerformTransaction(user_from, user_to, address, file_name)

	// wait :)
	<-done

	// delete transaction
	database.DeleteTransaction(address)

}
