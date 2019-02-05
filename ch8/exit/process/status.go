package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}

func main() {
	cmd := exec.Command("ls", "__a__")
	if err := cmd.Run(); err != nil {
		if status := exitStatus(cmd.ProcessState); status == -1 {
			fmt.Println(err)
		} else {
			fmt.Println("Status:", status)
		}
	}
}
