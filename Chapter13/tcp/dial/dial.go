package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify an address.")
	}
	addr := os.Args[1]
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalln("Listener:", addr, err)
		}
		time.Sleep(time.Second)
		c, err := listener.Accept()
		if err != nil {
			log.Fatalln("Listener:", addr, err)
		}
		log.Println("<- Connection to", addr)
		c.Close()
	}()
	ctx, canc := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer canc()
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", os.Args[1])
	if err != nil {
		log.Fatalln("-> Connection:", err)
	}
	log.Println("-> Connection to", os.Args[1])
	conn.Close()
}
