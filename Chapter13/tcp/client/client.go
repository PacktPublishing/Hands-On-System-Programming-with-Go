package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify an address.")
	}
	d := net.Dialer{}
	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	conn, err := d.DialContext(ctx, "tcp", os.Args[1])
	if err != nil {
		log.Fatalln("-> Connection:", err)
	}
	log.Println("-> Connection to", os.Args[1])
	var ch = make(chan []byte)
	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			fmt.Print("# ")
			msg, err := r.ReadBytes('\n')
			if err != nil {
				log.Println("-> Message error:", err)
			}
			if bytes.Compare(msg, []byte{'\\', 'q'}) == 0 {
				canc()
				continue
			}
			ch <- msg
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if _, err := conn.Write(msg); err != nil {
				log.Println("-> Connection:", err)
				return
			}
		}
	}
}
