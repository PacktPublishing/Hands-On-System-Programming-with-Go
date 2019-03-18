package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
)

type color int

func (c color) Start(w io.Writer) {
	fmt.Fprintf(w, "\x1b[%dm", c)
}

func (c color) End(w io.Writer) {
	fmt.Fprintf(w, "\x1b[%dm", 0)
}

func (c color) Fprintf(w io.Writer, format string, args ...interface{}) {
	c.Start(w)
	fmt.Fprintf(w, format, args...)
	c.End(w)
}

// List of colors
const (
	Red     color = 31
	Green   color = 32
	Yellow  color = 33
	Blue    color = 34
	Magenta color = 35
	Cyan    color = 36
	White   color = 37
)

func main() {
	w := os.Stdout
	for _, c := range []color{Red, Green, Yellow, Blue, Magenta, Cyan, White} {
		c.Fprintf(w, "the answer is %v\n", 42)
	}
	func(w io.Writer, args ...string) bool {
		rand.Shuffle(len(args), func(i, j int) {
			args[i], args[j] = args[j], args[i]
		})
		for i := range args {
			if i > 0 {
				fmt.Fprint(w, " ")
			}
			var f func(w io.Writer, format string, args ...interface{})
			if i%2 == 0 {
				f = Red.Fprintf
			} else {
				f = Green.Fprintf
			}
			f(w, "%s", args[i])
		}
		fmt.Fprintln(w)
		return false
	}(w, "banana", "pear", "apple", "something else")
}
