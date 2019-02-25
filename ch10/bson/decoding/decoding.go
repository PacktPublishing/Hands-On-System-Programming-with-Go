package main

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

type Character struct {
	Name        string `bson:"name"`
	Surname     string `bson:"surname"`
	Job         string `bson:"job,omitempty"`
	YearOfBirth int    `bson:"year_of_birth,omitempty"`
}

func main() {
	r := []byte(",\x00\x00\x00\x02name\x00\a\x00\x00" +
		"\x00Robert\x00\x02surname\x00\t\x00\x00\x00" +
		"Olmstead\x00\x00")
	var c Character
	if err := bson.Unmarshal(r, &c); err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", c)
}
