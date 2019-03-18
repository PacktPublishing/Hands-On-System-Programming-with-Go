package main

import (
	"log"

	. "github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch09/rpc/common"
)

func main() {
	r := ReadingList{}
	var books = []Book{
		{ISBN: "1", Title: "HELLO", Pages: 100},
		{ISBN: "2", Title: "JELLO", Pages: 50},
		{ISBN: "3", Title: "TRELLO", Pages: 150},
		{ISBN: "4", Title: "NELLO", Pages: 75},
	}
	for _, b := range books {
		if err := r.AddBook(b); err != nil {
			log.Fatal(err)
		}
	}
	if err := r.AddBook(books[0]); err != ErrDuplicate {
		log.Fatal(err)
	}
	if err := r.RemoveBook(books[1].ISBN); err != nil {
		log.Fatal(err)
	}
	if err := r.SetProgress(books[1].ISBN, 100); err != ErrMissing {
		log.Fatal(err)
	}
	if err := r.SetProgress(books[0].ISBN, 150); err != nil {
		log.Fatal(err)
	}
	if err := r.SetProgress(books[2].ISBN, 150); err != nil {
		log.Fatal(err)
	}
	if err := r.AdvanceProgress(books[3].ISBN, 200); err != nil {
		log.Fatal(err)
	}
	log.Println(r)

}
