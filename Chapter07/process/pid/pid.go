package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "-l")
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Cmd: ", cmd.Args[0])
	fmt.Println("Args:", cmd.Args[1:])
	fmt.Println("PID: ", cmd.Process.Pid)
	cmd.Wait()
}
