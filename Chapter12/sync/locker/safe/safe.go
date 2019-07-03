package main

import (
	"fmt"
	"sync"
)

func main() {
	done := make(chan struct{}, 10000)
	var a = 0
	m := sync.Mutex{}
	for i := 0; i < cap(done); i++ {
		go func(l sync.Locker, i int) {
			l.Lock()
			defer l.Unlock()
			if i%2 == 0 {
				a++
			} else {
				a--
			}
			done <- struct{}{}
		}(&m, i)
	}
	for i := 0; i < cap(done); i++ {
		<-done
	}
	fmt.Println(a)
}
