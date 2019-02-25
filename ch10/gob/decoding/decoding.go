package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Character struct {
	Name        string `gob:"name"`
	Surname     string `gob:"surname"`
	Job         string `gob:"job,omitempty"`
	YearOfBirth int    `gob:"year_of_birth,omitempty"`
}

func main() {
	data := []byte("D\xff\x81\x03\x01\x01\tCharacter" +
		"\x01\xff\x82\x00\x01\x04\x01\x04Name" +
		"\x01\f\x00\x01\aSurname\x01\f\x00\x01\x03" +
		"Job\x01\f\x00\x01\vYearOfBirth\x01\x04\x00" +
		"\x00\x00*\xff\x82\x01\x06Albert\x01\bWilmarth" +
		"\x01\x13assistant professor\x00")

	runDecode(data, new(Character))
	runDecode(data, new(struct {
		YearOfBirth int    `gob:"year_of_birth,omitempty"`
		Surname     string `gob:"surname"`
		Name        string `gob:"name"`
		Job         string `gob:"job,omitempty"`
	}))
	runDecode(data, new(struct {
		Name string `gob:"name"`
	}))
	runDecode(data, new(struct {
		Name        string `gob:"name"`
		Surname     string `gob:"surname"`
		Country     string `gob:"country"`
		Job         string `gob:"job,omitempty"`
		YearOfBirth int    `gob:"year_of_birth,omitempty"`
	}))
	runDecode(data, new(struct {
		Name []byte `gob:"name"`
	}))
}

func runDecode(data []byte, v interface{}) {
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(v); err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", v)
}
