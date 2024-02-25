package main

import (
	"Example/src/database"
	"Example/src/network"
	"os"
	"strconv"
)

// this is the daemon
func main() {
	args := os.Args[1:]
	user_to := args[0]
	file_name := args[1]

	node := network.Initialize_node()

	// I will computer and provide key
	key := network.Server(node)

	// initialize user
	database.HandleAccountStartup()
	user_from, _, _, _ := database.GetUserDetails()

	database.PerformTransaction(strconv.Itoa(user_from), user_to, key, file_name)

}
