package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	initdFile    = "/etc/init.d/mydaemon"
	varDir       = "/var/mydaemon/"
	pidFile      = "mydaemon.pid"
	outFile      = "mydaemon.log"
	errFile      = "mydaemon.err"
	initdContent = `#!/bin/sh

### BEGIN INIT INFO
# Provides:          mydaemon
# Required-Start:    $remote_fs $syslog
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Custom daemon
# Description:       Enable service provided by daemon.
### END INIT INFO

"%s" $1
`
)

// ErrSudo is an error that suggest to execute the command as super user
var ErrSudo error

var (
	bin string
	cmd string
)

func init() {
	p, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	bin = p
	if len(os.Args) != 1 {
		cmd = os.Args[1]
	}
	ErrSudo = fmt.Errorf("try `sudo %s %s`", bin, cmd)
}

func main() {
	var err error
	switch cmd {
	case "run":
		err = runApp()
	case "install":
		err = installApp()
	case "uninstall":
		err = uninstallApp()
	case "status":
		err = statusApp()
	case "start":
		err = startApp()
	case "stop":
		err = stopApp()
	default:
		helpApp()
	}
	if err != nil {
		fmt.Println(cmd, "error:", err)
	}
}

func installApp() error {
	_, err := os.Stat(initdFile)
	if err == nil {
		return errors.New("Already installed")
	}
	f, err := os.OpenFile(initdFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		if !os.IsPermission(err) {
			return err
		}
		return ErrSudo
	}
	defer f.Close()
	if _, err = fmt.Fprintf(f, initdContent, bin); err != nil {
		return err
	}
	fmt.Println("Daemon", bin, "installed")
	return nil
}

func uninstallApp() error {
	_, err := os.Stat(initdFile)
	if err != nil && os.IsNotExist(err) {
		return errors.New("not installed")
	}
	if err = os.Remove(initdFile); err != nil {
		if err != nil {
			if !os.IsPermission(err) {
				return err
			}
			return ErrSudo
		}
	}
	fmt.Println("Daemon", bin, "removed")
	return err
}

func statusApp() (err error) {
	var pid int
	defer func() {
		if pid == 0 {
			fmt.Println("status: not active")
			return
		}
		fmt.Println("status: active - pid", pid)
	}()
	pid, err = getPid()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	if err = p.Signal(syscall.Signal(0)); err != nil {
		fmt.Println(pid, "not found - removing PID file...")
		os.Remove(filepath.Join(varDir, pidFile))
		pid = 0
	}
	return nil
}

func writePid(pid int) (err error) {
	f, err := os.OpenFile(filepath.Join(varDir, pidFile), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = fmt.Fprintf(f, "%d", pid); err != nil {
		return err
	}
	return nil
}

func getPid() (pid int, err error) {
	b, err := ioutil.ReadFile(filepath.Join(varDir, pidFile))
	if err != nil {
		return 0, err
	}
	if pid, err = strconv.Atoi(string(b)); err != nil {
		return 0, fmt.Errorf("Invalid PID value: %s", string(b))
	}
	return pid, nil
}

func startApp() (err error) {
	const perm = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	if err = os.MkdirAll(varDir, 0755); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		return ErrSudo
	}
	cmd := exec.Command(bin, "run")
	cmd.Stdout, err = os.OpenFile(filepath.Join(varDir, outFile), perm, 0644)
	if err != nil {
		return err
	}
	cmd.Stderr, err = os.OpenFile(filepath.Join(varDir, errFile), perm, 0644)
	if err != nil {
		return err
	}
	cmd.Dir = "/"
	if err = cmd.Start(); err != nil {
		return err
	}
	if err := writePid(cmd.Process.Pid); err != nil {
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("Cannot kill process", cmd.Process.Pid, err)
		}
		return err
	}
	fmt.Println("Started with PID", cmd.Process.Pid)
	return nil
}

func stopApp() (err error) {
	pid, err := getPid()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	if err = p.Signal(os.Kill); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(varDir, pidFile)); err != nil {
		return err
	}
	fmt.Println("Stopped PID", pid)
	return nil
}

func runApp() error {
	fmt.Println("RUN")
	for {
		time.Sleep(time.Second)
	}
	return nil
}

func helpApp() error {
	fmt.Println("usage:", bin, "run|install|uninstall|status|start|stop")
	return nil
}
