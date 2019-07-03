package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 1; i <= 10; i++ {
		go func(a int) {
			for i := 1; i <= 10; i++ {
				fmt.Printf("%dx%d=%d\n", a, i, a*i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
