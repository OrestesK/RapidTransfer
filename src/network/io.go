package network

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Send_file(user_to string, filename string) {
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

func Receive_file(transaction_identifier string) {
	// node := Initialize_node()
	// TODO
	// database get key (big one) given transaction id (small one)
	// Client(node)
}
