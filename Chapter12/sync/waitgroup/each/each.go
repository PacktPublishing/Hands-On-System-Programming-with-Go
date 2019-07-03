package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 1; rand.Intn(10) != 0; i++ {
		wg.Add(1)
		go func(a int) {
			for i := 1; i <= 10; i++ {
				fmt.Printf("%dx%d=%d\n", a, i, a*i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
