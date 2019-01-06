package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/agnivade/levenshtein"
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

type argsScanner []string

func (a *argsScanner) Reset() { *a = (*a)[0:0] }

func (a *argsScanner) Parse(r io.Reader) (extra string) {
	s := bufio.NewScanner(r)
	s.Split(ScanArgs)
	for s.Scan() {
		*a = append(*a, s.Text())
	}
	if len(*a) == 0 {
		return ""
	}
	lastArg := (*a)[len(*a)-1]
	if !isQuote(rune(lastArg[0])) {
		return ""
	}
	*a = (*a)[:len(*a)-1]
	return lastArg + "\n"
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	w := os.Stdout
	a := argsScanner{}
	b := bytes.Buffer{}
	fmt.Fprint(w, "** Welcome to PseudoTerm! **\nPlease enter a command.\n")
	for {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot get working directory:", err)
			return
		}
		fmt.Fprintf(w, "\n[%s] > ", filepath.Base(pwd))

		a.Reset()
		b.Reset()
		for {
			s.Scan()
			b.Write(s.Bytes())
			extra := a.Parse(&b)
			if extra == "" {
				break
			}
			b.WriteString(extra)
		}

		idx := -1
		for i := range cmds {
			if !cmds[i].Match(a[0]) {
				continue
			}
			idx = i
			break
		}
		if idx == -1 {
			commandNotFound(w, a[0])
			continue
		}
		if cmds[idx].Run(w, a[1:]...) {
			fmt.Fprintln(w)
			return
		}
	}
}

func commandNotFound(w io.Writer, cmd string) {
	var list []string
	for _, c := range cmds {
		d := levenshtein.ComputeDistance(c.Name, cmd)
		if d < 3 {
			list = append(list, c.Name)
		}
	}
	fmt.Fprintf(w, "Command %q not found.", cmd)
	if len(list) == 0 {
		return
	}
	fmt.Fprint(w, " Maybe you meant: ")
	for i := range list {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprintf(w, "%s", list[i])
	}
}

var ErrClosingQuote = errors.New("Missing closing quote")

func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

func ScanArgs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start, first := 0, rune(0)
	for width := 0; start < len(data); start += width {
		first, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(first) {
			break
		}
	}
	if isQuote(first) {
		start++
	}
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if ok := isQuote(first); !ok && unicode.IsSpace(r) || ok && r == first {
			return i + width, data[start:i], nil
		}
	}

	if atEOF && len(data) > start {
		if isQuote(first) {
			start--
		}
		return len(data), data[start:], nil
	}
	if isQuote(first) {
		start--
	}
	return start, nil, nil
}
