package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var (
		v  atomic.Value
		wg sync.WaitGroup
	)
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("load", v.Load())
			wg.Done()
		}(i)
		go func(i int) {
			v.Store(i)
			fmt.Println("store", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
