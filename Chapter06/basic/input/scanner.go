package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

type cmd struct {
	Name   string
	Help   string
	Action func(w io.Writer, args ...string) bool
}

func (c cmd) Match(s string) bool {
	return c.Name == s
}

func (c cmd) Run(w io.Writer, args ...string) bool {
	return c.Action(w, args...)
}

var cmds = make([]cmd, 0, 10)

func init() {
	cmds = append(cmds,
		cmd{
			Name: "exit",
			Help: "Exits the application",
			Action: func(w io.Writer, args ...string) bool {
				fmt.Fprintf(w, "Goodbye! :)\n")
				return true
			},
		},
		cmd{
			Name: "help",
			Help: "Shows available commands",
			Action: func(w io.Writer, args ...string) bool {
				fmt.Fprintln(w, "Available commands:")
				for _, c := range cmds {
					fmt.Fprintf(w, "  - %-15s %s\n", c.Name, c.Help)
				}
				return false
			},
		},
		cmd{
			Name: "shuffle",
			Help: "Shuffles a list of strings",
			Action: func(w io.Writer, args ...string) bool {
				rand.Shuffle(len(args), func(i, j int) {
					args[i], args[j] = args[j], args[i]
				})
				for i := range args {
					if i > 0 {
						fmt.Fprint(w, " ")
					}
					fmt.Fprintf(w, "%s", args[i])
				}
				fmt.Fprintln(w)
				return false
			},
		},
		cmd{
			Name: "print",
			Help: "Prints a file",
			Action: func(w io.Writer, args ...string) bool {
				if len(args) != 1 {
					fmt.Fprintln(w, "Please specify one file!")
					return false
				}
				f, err := os.Open(args[0])
				if err != nil {
					fmt.Fprintf(w, "Cannot open %s: %s\n", args[0], err)
				}
				defer f.Close()
				if _, err := io.Copy(w, f); err != nil {
					fmt.Fprintf(w, "Cannot print %s: %s\n", args[0], err)
				}
				fmt.Fprintln(w)
				return false
			},
		},
	)
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	w := os.Stdout
	fmt.Fprint(w, "** Welcome to PseudoTerm! **\nPlease enter a command.\n")
	for {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot get working directory:", err)
			return
		}
		fmt.Fprintf(w, "\n[%s] > ", filepath.Base(pwd))
		if !s.Scan() {
			continue
		}
		args := strings.Split(string(s.Bytes()), " ")
		idx := -1
		for i := range cmds {
			if !cmds[i].Match(args[0]) {
				continue
			}
			idx = i
			break
		}
		if idx == -1 {
			fmt.Fprintf(w, "%q not found. Use `help` for available commands\n", args[0])
			continue
		}
		if cmds[idx].Run(w, args[1:]...) {
			fmt.Fprintln(w)
			return
		}
	}
}
