package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Greeter interface {
	Greet(w io.Writer)
}

type Character struct {
	Name        string `gob:"name"`
	Surname     string `gob:"surname"`
	Job         string `gob:"job,omitempty"`
	YearOfBirth int    `gob:"year_of_birth,omitempty"`
}

func (c Character) Greet(w io.Writer) {
	fmt.Fprintf(w, "Hello, my name is %s %s", c.Name, c.Surname)
	if c.Job != "" {
		fmt.Fprintf(w, " and I am a %s", c.Job)
	}
}

func main() {
	gob.Register(Greeter(Character{}))
	r := strings.NewReader("U\x10\x00\x0emain.Character" +
		"\xff\x81\x03\x01\x01\tCharacter\x01\xff\x82\x00" +
		"\x01\x04\x01\x04Name\x01\f\x00\x01\aSurname" +
		"\x01\f\x00\x01\x03Job\x01\f\x00\x01\vYearOfBirth" +
		"\x01\x04\x00\x00\x00\x1f\xff\x82\x1c\x01\x05John" +
		" \x01\aKirowan\x01\tprofessor\x00")
	var char Greeter
	if err := gob.NewDecoder(r).Decode(&char); err != nil {
		log.Fatalln(err)
	}
	char.Greet(os.Stdout)
}
