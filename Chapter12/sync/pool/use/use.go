package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

func Get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

func Put(b *bytes.Buffer) {
	b.Reset()
	pool.Put(b)
}

func main() {
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(20)
	var (
		list = make([][]byte, 20)
		m    sync.Mutex
	)
	for i := 0; i < 20; i++ {
		go func(v int) {
			time.Sleep(time.Second * time.Duration(1+v/4))
			b := Get()
			defer func() {
				Put(b)
				wg.Done()
			}()
			fmt.Fprintf(b, "Goroutine %2d using %p, after %.0fs\n", v, b, time.Since(start).Seconds())
			m.Lock()
			list[v] = b.Bytes()
			m.Unlock()
		}(i)
	}
	wg.Wait()
	for i := range list {
		fmt.Printf("%d - %s", i, list[i])
	}
}
