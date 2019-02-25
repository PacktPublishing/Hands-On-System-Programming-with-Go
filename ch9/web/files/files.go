package main

import (
	"errors"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify a directory")
	}
	s, err := os.Stat(os.Args[1])
	if err == nil && !s.IsDir() {
		err = errors.New("not a directory")
	}
	if err != nil {
		log.Fatalln("Invalid path:", err)
	}
	http.Handle("/", http.FileServer(http.Dir(os.Args[1])))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
