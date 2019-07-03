package main

import (
	"log"
	"strings"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch15/use/prop"
)

type UpperString string

func (u *UpperString) UnmarshalProp(b []byte) error {
	*u = UpperString(strings.ToUpper(string(b)))
	return nil
}

func main() {
	r := strings.NewReader(
		"\n# comment, ignore\nkey1: 10.5\nkey2: some string" +
			"\nkey3: 42\nkey4: false\nspecial: another string\n")
	var v struct {
		Key1 float32
		Key2 string
		Key3 uint64
		Key4 bool
		Key5 UpperString `prop:"special"`
		key6 int
	}
	if err := prop.NewDecoder(r).Decode(&v); err != nil {
		log.Fatal(r)
	}
	log.Printf("%+v", v)
}
