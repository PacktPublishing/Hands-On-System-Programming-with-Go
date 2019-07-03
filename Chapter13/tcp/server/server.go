package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify an address.")
	}
	addr := os.Args[1]

	ctx, canc := context.WithCancel(context.Background())
	defer canc()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		canc()
	}()

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
	if err != nil {
		log.Fatalln("Listener:", addr, err)
	}
	var conn = make(chan net.Conn)
	go func() {
		log.Println("start...")
		for {
			c, err := listener.Accept()
			if err != nil {
				log.Fatalln("<- Accept:", addr, err)
			}
			conn <- c
		}
	}()
loop:
	for {
		select {
		case c := <-conn:
			go handleConn(ctx, canc, c)
		case <-ctx.Done():
			break loop
		}
	}
	log.Println("exit...")
}

func handleConn(ctx context.Context, canc func(), conn net.Conn) {
	ch := make(chan string)
	defer conn.Close()
	go func() {
		r := bufio.NewReader(conn)
		defer close(ch)
		for {
			msg, err := r.ReadString('\n')
			if err != nil {
				if nerr, ok := err.(net.Error); ok && !nerr.Temporary() {
					log.Println("<- Network error:", err)
					return
				}
				if err == io.EOF {
					return
				}
				log.Println("<- Message error:", err)
				continue
			}
			select {
			case ch <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			switch msg = strings.TrimSpace(msg); msg {
			case `\q`:
				log.Println("<- Exiting...")
				canc()
				return
			case `\x`:
				log.Println("<- Special message `\\x` received!")
			default:
				log.Println("<- Message Received:", msg)
			}
		}
	}
}
