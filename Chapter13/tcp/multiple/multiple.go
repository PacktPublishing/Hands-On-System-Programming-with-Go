package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	list := []string{
		"localhost:9090",
		"localhost:9091",
		"localhost:9092",
	}
	go func() {
		listener, err := net.Listen("tcp", list[0])
		if err != nil {
			log.Fatalln("Listener:", list[0], err)
		}
		time.Sleep(time.Second * 5)
		c, err := listener.Accept()
		if err != nil {
			log.Fatalln("Listener:", list[0], err)
		}
		defer c.Close()
	}()
	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	wg := sync.WaitGroup{}
	wg.Add(len(list))
	for _, addr := range list {
		go func(addr string) {
			defer wg.Done()
			conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
			if err != nil {
				log.Println("-> Connection:", err)
				return
			}
			log.Println("-> Connection to", addr, "cancelling context")
			canc()
			conn.Close()
		}(addr)
	}
	wg.Wait()
}
