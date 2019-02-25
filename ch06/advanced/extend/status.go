package main

import (
	"fmt"
	"io"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch06/advanced/extend/command"
)

func init() {
	command.Register(&Stack{})
}

type Stack struct{ data []string }

func (s *Stack) push(values ...string) {
	s.data = append(s.data, values...)
}

func (s *Stack) pop() (string, bool) {
	if len(s.data) == 0 {
		return "", false
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v, true
}

func (s *Stack) GetName() string {
	return "stack"
}

func (s *Stack) GetHelp() string {
	return "a stack-like memory storage"
}

func (s *Stack) isValid(cmd string, args []string) bool {
	switch cmd {
	case "pop":
		return len(args) == 0
	case "push":
		return len(args) > 0
	default:
		return false
	}
}

func (s *Stack) Run(r io.Reader, w io.Writer, args ...string) (exit bool) {
	if l := len(args); l < 2 || !s.isValid(args[1], args[2:]) {
		fmt.Fprintf(w, "Use `stack push <something>` or `stack pop`\n")
		return false
	}
	if args[1] == "push" {
		s.push(args[2:]...)
		return false
	}
	if v, ok := s.pop(); !ok {
		fmt.Fprintf(w, "Empty!\n")
	} else {
		fmt.Fprintf(w, "Got: `%s`\n", v)
	}
	return false
}
