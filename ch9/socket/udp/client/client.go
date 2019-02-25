package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify an address.")
	}
	addr, err := net.ResolveUDPAddr("udp", os.Args[1])
	if err != nil {
		log.Fatalln("Invalid address:", os.Args[1], err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalln("-> Connection:", err)
	}
	log.Println("-> Connection to", addr)
	r := bufio.NewReader(os.Stdin)
	b := make([]byte, 1024)
	for {
		fmt.Print("# ")
		msg, err := r.ReadBytes('\n')
		if err != nil {
			log.Println("-> Message error:", err)
		}
		if _, err := conn.Write(msg); err != nil {
			log.Println("-> Connection:", err)
			return
		}
		n, err := conn.Read(b)
		if err != nil {
			log.Println("<- Receive error:", err)
		}
		msg = bytes.TrimSpace(b[:n])
		log.Printf("<- %q", msg)
	}
}
