package main

import (
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type Character struct {
	Name        string `yaml:"name"`
	Surname     string `yaml:"surname"`
	Job         string `yaml:"job,omitempty"`
	YearOfBirth int    `yaml:"year_of_birth,omitempty"`
}

func main() {
	r := strings.NewReader(`- name: John Raymond
  surname: Legrasse
  job: policeman
- name: "Francis"
  surname: Wayland Thurston
  job: anthropologist`)
	d := yaml.NewDecoder(r)
	var c []Character
	if err := d.Decode(&c); err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", c)
}
