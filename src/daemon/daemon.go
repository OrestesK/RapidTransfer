package main

import (
	"Rapid/src/database"
	"Rapid/src/network"
	"os"
)

// this is the daemon
func main() {
	args := os.Args[1:]
	user_from := args[0]
	user_to := args[1]
	file_name := args[2]

	//fmt.Printf("to user: %s filename: %s\n", user_to, file_name)

	node := network.Initialize_node()

	done := make(chan bool)
	// I will computer and provide key
	address := network.Server(node, file_name, done)

	database.InitializeDatabase()
	database.PerformTransaction(user_from, user_to, address, file_name)
	println("Waiting for user to accept file")

	// wait :)
	<-done

	// delete transaction
	database.DeleteTransactionWithAddress(address)
}
