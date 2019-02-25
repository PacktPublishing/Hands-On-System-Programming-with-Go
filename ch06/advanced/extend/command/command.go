package command

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/agnivade/levenshtein"
)

// ErrDuplicateCommand is returned when two commands have the same name
var ErrDuplicateCommand = errors.New("Duplicate command")

// Command represents a terminal command
type Command interface {
	GetName() string
	GetHelp() string
	Run(input io.Reader, output io.Writer, args ...string) (exit bool)
}

var commands []Command

// Register adds the Command to the command list
func Register(command Command) error {
	name := command.GetName()
	for i, c := range commands {
		// unique commands in alphabetical order
		switch strings.Compare(c.GetName(), name) {
		case 0:
			return ErrDuplicateCommand
		case 1:
			commands = append(commands, nil)
			copy(commands[i+1:], commands[i:])
			commands[i] = command
			return nil
		case -1:
			continue
		}
	}
	commands = append(commands, command)
	return nil
}

// GetCommand returns the command with the given name
func GetCommand(name string) Command {
	for _, c := range commands {
		if c.GetName() == name {
			return c
		}
	}
	return suggest
}

// Base is a basic Command that runs a closure
type Base struct {
	Name, Help string
	Action     func(input io.Reader, output io.Writer, args ...string) (exit bool)
}

func (b Base) String() string { return b.Name }

// GetName returns the Name
func (b Base) GetName() string { return b.Name }

// GetHelp returns the Help
func (b Base) GetHelp() string { return b.Help }

// Run calls the closure
func (b Base) Run(input io.Reader, output io.Writer, args ...string) (exit bool) {
	return b.Action(input, output, args...)
}

func init() {
	Register(Base{Name: "help", Help: "Shows available commands", Action: helpAction})
	Register(Base{Name: "exit", Help: "Exits the application", Action: exitAction})
}

func helpAction(in io.Reader, w io.Writer, args ...string) bool {
	fmt.Fprintln(w, "Available commands:")
	for _, c := range commands {
		n := c.GetName()
		fmt.Fprintf(w, "  - %-15s %s\n", n, c.GetHelp())
	}
	return false
}

func exitAction(in io.Reader, w io.Writer, args ...string) bool {
	fmt.Fprintf(w, "Goodbye! :)\n")
	return true
}

var suggest = Base{
	Action: func(in io.Reader, w io.Writer, args ...string) bool {
		var list []string
		for _, c := range commands {
			name := c.GetName()
			d := levenshtein.ComputeDistance(name, args[0])
			if d < 3 {
				list = append(list, name)
			}
		}
		fmt.Fprintf(w, "Command %q not found.", args[0])
		if len(list) == 0 {
			return false
		}
		fmt.Fprint(w, " Maybe you meant: ")
		for i := range list {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}
			fmt.Fprintf(w, "%s", list[i])
		}
		return false
	},
}
