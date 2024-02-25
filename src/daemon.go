package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func execute_daemon() {
	// execute daemon, runs in background independent of this
	cmd := exec.Command("go", "run", "src/daemon.go")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Start()
	syscall.Setpgid(cmd.Process.Pid, cmd.Process.Pid)

	// saved pid to file
	tt := fmt.Sprintf("%d\n", cmd.Process.Pid)
	os.WriteFile("pid", []byte(tt), 0755)
}
