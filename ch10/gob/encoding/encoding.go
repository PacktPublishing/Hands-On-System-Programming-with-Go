package main

import (
	"encoding/gob"
	"log"
	"strings"
)

type Character struct {
	Name        string `gob:"name"`
	Surname     string `gob:"surname"`
	Job         string `gob:"job,omitempty"`
	YearOfBirth int    `gob:"year_of_birth,omitempty"`
}

func main() {
	var char = Character{
		Name:    "Albert",
		Surname: "Wilmarth",
		Job:     "assistant professor",
	}
	s := strings.Builder{}
	e := gob.NewEncoder(&s)
	if err := e.Encode(char); err != nil {
		log.Fatalln(err)
	}
	log.Printf("%q", s.String())
}
