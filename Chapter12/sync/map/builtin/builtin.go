package main

import (
	"fmt"
	"sync"
)

func main() {
	var m = sync.Map{}
	var wg = sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			m.LoadOrStore(i, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	i := 0
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	fmt.Println(i)
}
