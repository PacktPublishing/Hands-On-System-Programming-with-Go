package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	f, err := os.OpenFile("file.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	go func() {
		<-c
		w.Flush()
		os.Exit(0)
	}()
	for i := 0; i < 3; i++ {
		fmt.Fprintln(w, "hello")
		log.Println(i)
		time.Sleep(time.Second)
	}
}
