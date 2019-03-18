package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch09/rpc/common"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify an address.")
	}
	if err := rpc.Register(&common.ReadingService{}); err != nil {
		log.Fatalln(err)
	}
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Server Started")
	if err := http.Serve(l, nil); err != nil {
		log.Fatal(err)
	}
}
