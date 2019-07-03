package main

import (
	"bytes"
	"fmt"
	"os"
)

var (
	b   = make([]byte, 128)
	buf = bytes.NewBuffer(b)
)

const grr = "G.R.R. Martin"

type book struct {
	Author, Title string
	Year          int
}

func (b book) String() {
	buf.Reset()
}

func main() {
	dst := os.Stdout
	defer dst.Close()
	bookList := []book{
		{Author: grr, Title: "A Game of Thrones", Year: 1996},
		{Author: grr, Title: "A Clash of Kings", Year: 1998},
		{Author: grr, Title: "A Storm of Swords", Year: 2000},
		{Author: grr, Title: "A Feast for Crows", Year: 2005},
		{Author: grr, Title: "A Dance with Dragons", Year: 2011},
		{Author: grr, Title: "The Winds of Winter"},
		{Author: grr, Title: "A Dream of Spring"},
	}
	b := bytes.NewBuffer(make([]byte, 0, 16))
	for _, v := range bookList {
		// prints a msg formatted with arguments to writer
		fmt.Fprintf(b, "%s - %s", v.Title, v.Author)
		if v.Year > 0 { // we do not print the year if it's not there
			fmt.Fprintf(b, " (%d)", v.Year)
		}
		b.WriteRune('\n')
		if _, err := b.WriteTo(dst); err != nil { // copies bytes, drains buffer
			fmt.Println("Error:", err)
			return
		}
	}
}
