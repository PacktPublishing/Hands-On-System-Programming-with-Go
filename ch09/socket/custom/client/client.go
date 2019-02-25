package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch09/socket/custom/common"
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
		data, err := common.CreateMessage(bytes.TrimSpace(msg))
		if err != nil {
			log.Println("-> Encode error:", err)
		}
		if _, err := conn.Write(data); err != nil {
			log.Println("-> Connection:", err)
			return
		}
		n, err := conn.Read(b)
		if err != nil {
			log.Println("<- Receive error:", err)
			continue
		}
		msg, err = common.MessageContent(b[:n])
		if err != nil {
			log.Println("<- Decode error:", err)
			continue
		}
		log.Printf("<- %q", msg)
	}
}
