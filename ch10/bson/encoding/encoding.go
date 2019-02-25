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
	var char = Character{
		Name:    "Robert",
		Surname: "Olmstead",
	}
	b, err := bson.Marshal(char)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%q", b)
}
