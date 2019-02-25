package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch06/advanced/state/command"
)

func init() {
	command.Register(command.Base{
		Name:   "shuffle",
		Help:   "Shuffles a list of strings",
		Action: shuffleAction,
	})
	command.Register(command.Base{
		Name:   "print",
		Help:   "Prints a file",
		Action: printAction,
	})
}

func main() {
	s, w, a, b := bufio.NewScanner(os.Stdin), os.Stdout, args{}, bytes.Buffer{}
	command.Startup(w)
	defer command.Shutdown(w)
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
		if command.GetCommand(a[0]).Run(os.Stdin, w, a...) {
			fmt.Fprintln(w)
			return
		}
	}
}

type args []string

func (a *args) Reset() { *a = (*a)[0:0] }

func (a *args) Parse(r io.Reader) (extra string) {
	s := bufio.NewScanner(r)
	s.Split(scanargs)
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

func isQuote(r rune) bool { return r == '"' || r == '\'' }

func scanargs(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

func shuffleAction(r io.Reader, w io.Writer, args ...string) bool {
	args = args[1:]
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
}

func printAction(r io.Reader, w io.Writer, args ...string) bool {
	if len(args) != 2 {
		fmt.Fprintln(w, "Please specify one file!")
		return false
	}
	f, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(w, "Cannot open %s: %s\n", args[1], err)
	}
	defer f.Close()
	if _, err := io.Copy(w, f); err != nil {
		fmt.Fprintf(w, "Cannot print %s: %s\n", args[1], err)
	}
	fmt.Fprintln(w)
	return false
}
