package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("Start application...")
	c := make(chan os.Signal)
	signal.Notify(c)
	log.Println("Exit with signal:", <-c)
}
