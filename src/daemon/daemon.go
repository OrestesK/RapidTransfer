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

	node := network.Initialize_node()

	done := make(chan bool)
	// I will computer and provide key
	address := network.Server(node, file_name, done)

	// initialize user
	database.HandleAccountStartup()
	_, user_from, _, _ := database.GetUserDetails()

	thing := database.PerformTransaction(user_from, user_to, address, file_name)
	println(thing)

	// wait :)
	<-done

	// delete transaction
	database.DeleteTransactionWithAddress(address)
}
