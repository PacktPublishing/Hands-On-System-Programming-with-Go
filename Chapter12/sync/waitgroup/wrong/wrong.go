package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		go func(a int) {
			wg.Add(1)
			for i := 1; i <= 10; i++ {
				fmt.Printf("%dx%d=%d\n", a, i, a*i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
