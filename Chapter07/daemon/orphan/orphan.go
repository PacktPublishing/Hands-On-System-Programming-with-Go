package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var pid = os.Getpid()

func main() {
	fmt.Printf("[%d] Start\n", pid)
	fmt.Printf("[%d] PPID: %d\n", pid, os.Getppid())
	defer fmt.Printf("[%d] Exit\n\n", pid)
	if len(os.Args) != 1 {
		runDaemon()
		return
	}
	if err := forkProcess(); err != nil {
		fmt.Printf("[%d] Fork error: %s\n", pid, err)
		return
	}
	if err := releaseResources(); err != nil {
		fmt.Printf("[%d] Release error: %s\n", pid, err)
		return
	}
}

func forkProcess() error {
	cmd := exec.Command(os.Args[0], "daemon")
	cmd.Stdout, cmd.Stderr, cmd.Dir = os.Stdout, os.Stderr, "/"
	return cmd.Start()
}

func releaseResources() error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Release()
}

func runDaemon() {
	for {
		fmt.Printf("[%d] Daemon mode\n", pid)
		time.Sleep(time.Second * 10)
	}
}
