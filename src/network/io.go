package network

import (
	"Example/src/database"
	"fmt"

	// "fmt"
	// "os"
	"os/exec"
)

func Send_file(user_to string, filename string) {
	// execute daemon, runs in background independent of this
	cmd := exec.Command("go", "run", "src/daemon/daemon.go", user_to, filename)

	cmd.Start()

	// saved pid to file
	// tt := fmt.Sprintf("%d\n", cmd.Process.Pid)
	// os.WriteFile("pid", []byte(tt), 0755)
}

func Receive_file(transaction_identifier string) {
	node := Initialize_node()
	_, name, _, _ := database.GetUserDetails()
	id := database.GetUserID(name)
	result := database.UserCanViewTransaction(id, transaction_identifier)
	if !result {
		// Cannot view this transaction.
		fmt.Print("You cannot download this file.")
		return
	}

	// get big key from small key
	address := database.GetAddressFromTransactionPhrase(transaction_identifier)
	file_name := database.GetFileNameFromTransactionPhrase(transaction_identifier)
	done := make(chan bool)
	// client
	if len(file_name) == 0 {
		fmt.Println("File Not found")
	} else {
		Client(node, address, file_name, done, false)
		<-done
	}
}

func Fake_receive_file(transaction_identifier string) {
	node := Initialize_node()

	// get big key from small key
	address := database.GetAddressFromTransactionPhrase(transaction_identifier)
	file_name := database.GetFileNameFromTransactionPhrase(transaction_identifier)

	done := make(chan bool)
	// client
	Client(node, address, file_name, done, true)
	<-done
}
