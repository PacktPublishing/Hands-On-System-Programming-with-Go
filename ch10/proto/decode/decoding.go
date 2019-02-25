package main

import (
	"log"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch10/proto/gen"
	"github.com/golang/protobuf/proto"
)

func main() {
	b := proto.NewBuffer([]byte(
		"/\n\x06George\x12\x0eGammell Angell" +
			"\x1a\x12professor emeritus \xaa\x0e",
	))
	var char gen.Character
	if err := b.DecodeMessage(&char); err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", char)
}
