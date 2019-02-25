package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func processState(e error) *os.ProcessState {
	err, ok := e.(*exec.ExitError)
	if !ok {
		return nil
	}
	return err.ProcessState
}

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
		state := processState(err)
		if state == nil {
			fmt.Println(err)
			return
		}
		if status := exitStatus(state); status == -1 {
			fmt.Println(err)
		} else {
			fmt.Println("Status:", status)
		}
	}
}
