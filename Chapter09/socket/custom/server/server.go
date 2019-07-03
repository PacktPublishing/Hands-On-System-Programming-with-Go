package main

import (
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
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln("Listener:", os.Args[1], err)
	}

	b := make([]byte, 256*256)
	for {
		n, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Println("<-", addr, "Message error:", err)
			continue
		}
		msg, err := common.MessageContent(b[:n])
		if err != nil {
			log.Println("<-", addr, "Decode error:", err)
			continue
		}
		log.Printf("<- %q from %s", msg, addr)
		for i, l := 0, len(msg); i < l/2; i++ {
			msg[i], msg[l-1-i] = msg[l-1-i], msg[i]
		}
		data, err := common.CreateMessage(msg)
		if err != nil {
			log.Println("->", addr, "Encode error:", err)
			continue
		}
		if _, err := conn.WriteTo(data, addr); err != nil {
			log.Println("->", addr, "Send error:", err)
		}
	}
}
