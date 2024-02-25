package main

import (
	"Example/src/database"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func execute_daemon() {
	// execute daemon, runs in background independent of this
	cmd := exec.Command("go", "run", "src/daemon.go, ")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Start()
	syscall.Setpgid(cmd.Process.Pid, cmd.Process.Pid)

	// saved pid to file
	tt := fmt.Sprintf("%d\n", cmd.Process.Pid)
	os.WriteFile("pid", []byte(tt), 0755)
}

// this is the daemon
func main() {
	args := os.Args[1:]
	user_from := args[0]
	user_to := args[1]
	file_name := args[2]
	// I will computer and provide key

	// TODO INITIALIZE DATABASE
	node := initialize_node()
	key := server(node)

	database.PerformTransaction(user_from, user_to, file_name, key)

}
