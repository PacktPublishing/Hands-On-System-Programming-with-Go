package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type record struct {
	sync.Mutex
	buf     string
	cond    *sync.Cond
	writers []io.Writer
}

func (r *record) Run() {
	for i := range r.writers {
		go func(i int) {
			for {
				r.Lock()
				r.cond.Wait()
				fmt.Fprintf(r.writers[i], "%s\n", r.buf)
				r.Unlock()
			}
		}(i)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify at least a file")
	}
	r := record{
		writers: make([]io.Writer, len(os.Args)-1),
	}
	r.cond = sync.NewCond(&r)
	for i, v := range os.Args[1:] {
		f, err := os.Create(v)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r.writers[i] = f
	}
	r.Run()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(":> ")
		if !scanner.Scan() {
			break
		}
		r.Lock()
		r.buf = scanner.Text()
		r.Unlock()
		switch {
		case r.buf == `\q`:
			return
		default:
			r.cond.Broadcast()
		}
	}
}
